package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/mailer"
	"github.com/gophish/gophish/models"
	"github.com/sirupsen/logrus"
)

// Worker is an interface that defines the operations needed for a background worker
type Worker interface {
	Start()
	LaunchCampaign(c models.Campaign)
	SendTestEmail(s *models.EmailRequest) error
	SendTestWhatsapp(s *models.EmailRequest) error
}

// DefaultWorker is the background worker that handles watching for new campaigns and sending emails appropriately.
type DefaultWorker struct {
	mailer            mailer.Mailer
	whatsappMessenger mailer.WhatsappMessenger
}

// New creates a new worker object to handle the creation of campaigns
func New(options ...func(Worker) error) (Worker, error) {
	log.Info("Creating new Worker")
	defaultMailer := mailer.NewMailWorker()
	defaultWhatsappSender := mailer.NewMessageWorker()

	w := &DefaultWorker{
		whatsappMessenger: defaultWhatsappSender,
		mailer:            defaultMailer,
	}
	for _, opt := range options {
		if err := opt(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

// WithMailer sets the mailer for a given worker.
// By default, workers use a standard, default mailworker.
func WithMailer(m mailer.Mailer) func(*DefaultWorker) error {
	return func(w *DefaultWorker) error {
		w.mailer = m
		return nil
	}
}

// WithWhatsappMessenger sets the whatsappMessenger for a given worker.
// By default, workers use a standard, default messageworker.
func WithWhatsappMessenger(ws mailer.WhatsappMessenger) func(*DefaultWorker) error {
	return func(w *DefaultWorker) error {
		w.whatsappMessenger = ws
		return nil
	}
}

// processCampaigns loads maillogs scheduled to be sent before the provided
// time and sends them to the mailer.
func (w *DefaultWorker) processCampaigns(t time.Time) error {
	log.Info("Processing campaigns")
	ms, err := models.GetQueuedMailLogs(t.UTC())
	if err != nil {
		log.Error(err)
		return err
	}
	// Lock the MailLogs (they will be unlocked after processing)
	err = models.LockMailLogs(ms, true)
	if err != nil {
		return err
	}
	campaignCache := make(map[int64]models.Campaign)
	// We'll group the maillogs by campaign ID to (roughly) group
	// them by sending profile. This lets the mailer re-use the Sender
	// instead of having to re-connect to the SMTP server for every
	// email.
	msg := make(map[int64][]mailer.Mail)
	wspMsg := make(map[int64][]mailer.Message)
	for _, m := range ms {
		// We cache the campaign here to greatly reduce the time it takes to
		// generate the message (ref #1726)
		log.Info(fmt.Sprintf("Caching campaign with ID %d", m.CampaignId))
		c, ok := campaignCache[m.CampaignId]
		if !ok {
			c, err = models.GetCampaignMailContext(m.CampaignId, m.UserId)
			if err != nil {
				return err
			}
			campaignCache[c.Id] = c
		}
		m.CacheCampaign(&c)
		if c.SMTP.Interface == "SMTP" {
			msg[m.CampaignId] = append(msg[m.CampaignId], m)
		} else if c.SMTP.Interface == "Whatsapp" {
			wspMsg[m.CampaignId] = append(wspMsg[m.CampaignId], m)
		}
	}

	// Next, we process each group of maillogs in parallel
	for cid, msc := range msg {
		go w.processMailGroup(campaignCache[cid], msc)
	}
	for cid, wsc := range wspMsg {
		go w.processWhatsappGroup(campaignCache[cid], wsc)
	}
	return nil
}

// setToInProgress sets campaign's status to in progress
func setToInProgress(c models.Campaign) {
	log.Info(fmt.Sprintf("Setting campaign %d to In Progress", c.Id))
	if c.Status == models.CampaignQueued {
		err := c.UpdateStatus(models.CampaignInProgress)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

// processMailGroup processes a group of maillogs
func (w *DefaultWorker) processMailGroup(c models.Campaign, msc []mailer.Mail) {
	log.Info(fmt.Sprintf("Processing SMTP campaign with ID %d", c.Id))
	setToInProgress(c)
	log.WithFields(logrus.Fields{
		"num_emails": len(msc),
	}).Info("Sending emails to mailer for processing")
	w.mailer.Queue(msc)
}

// processWhatsappGroup process a group of whatsapp messages
func (w *DefaultWorker) processWhatsappGroup(c models.Campaign, wspMsc []mailer.Message) {
	log.Info(fmt.Sprintf("Processing Whatsapp campaign with ID %d", c.Id))
	setToInProgress(c)
	log.WithFields(logrus.Fields{
		"num_whatsapps": len(wspMsc),
	}).Info("Sending whatsapps to messenger for processing")
	w.whatsappMessenger.MessageQueue(wspMsc)
}

// Start launches the worker to poll the database every minute for any pending maillogs
// that need to be processed.
func (w *DefaultWorker) Start() {
	log.Info("Background Worker Started Successfully - Waiting for Campaigns")
	go w.mailer.Start(context.Background())
	go w.whatsappMessenger.StartMessaging(context.Background())
	for t := range time.Tick(1 * time.Minute) {
		err := w.processCampaigns(t)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

// LaunchCampaign starts a campaign
func (w *DefaultWorker) LaunchCampaign(c models.Campaign) {
	log.Info(fmt.Sprintf("Launching campaign %d", c.Id))
	ms, err := models.GetMailLogsByCampaign(c.Id)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(fmt.Sprintf("Locking mails for campaign %d", c.Id))
	models.LockMailLogs(ms, true)
	// This is required since you cannot pass a slice of values
	// that implements an interface as a slice of that interface.
	mailEntries := []mailer.Mail{}
	whatsappEntries := []mailer.Message{}
	currentTime := time.Now().UTC()
	log.Info(fmt.Sprintf("Getting context for campaign %d", c.Id))
	campaignMailCtx, err := models.GetCampaignMailContext(c.Id, c.UserId)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(fmt.Sprintf("Sending scheduled messages for campaign %d", c.Id))
	for _, m := range ms {
		// Only send the emails scheduled to be sent for the past minute to
		// respect the campaign scheduling options
		if m.SendDate.After(currentTime) {
			m.Unlock()
			continue
		}
		err = m.CacheCampaign(&campaignMailCtx)
		if err != nil {
			log.Error(err)
			return
		}
		mailEntries = append(mailEntries, m)
		whatsappEntries = append(whatsappEntries, m)
	}
	log.Info(fmt.Sprintf("Queuing message entries for campaign %d", c.Id))
	switch c.SMTP.Interface {
	case "SMTP":
		w.mailer.Queue(mailEntries)
	case "Whatsapp":
		w.whatsappMessenger.MessageQueue(whatsappEntries)
	}
}

// SendTestEmail sends a test email
func (w *DefaultWorker) SendTestEmail(s *models.EmailRequest) error {
	go func() {
		ms := []mailer.Mail{s}
		w.mailer.Queue(ms)
	}()
	return <-s.ErrorChan
}

// SendTestWhatsapp sends a test Whatsapp message
func (w *DefaultWorker) SendTestWhatsapp(s *models.EmailRequest) error {
	// go func() {
	// 	ws := []mailer.Message{s}
	// 	w.whatsappMessenger.Queue(ws)
	// }()
	// return <- s.ErrorChan
	return errors.New("Method not implemented")
}
