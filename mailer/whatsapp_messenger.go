package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/gophish/gophish/logger"
	"github.com/sirupsen/logrus"
)

// WhatsappAPI is the URL to which messages are submitted
var WhatsappAPI = "https://graph.facebook.com/v21.0/%s/messages"

// WhatsappMaxReconnectAttempts is the maximum number of times we should reconnect to a server
var WhatsappMaxReconnectAttempts = 10

// ErrWhatsappMaxConnectAttempts is thrown when the maximum number of reconnect attempts
// is reached.
type ErrWhatsappMaxConnectAttempts struct {
	underlyingError error
}

// ErrWhatsappAPIResponse is thrown when Whatsapp API returns an error
var ErrWhatsappAPIResponse = errors.New("Whatsapp API response error")

// WhatsappError returns the wrapped error response
func (e *ErrMaxConnectAttempts) WhatsappError() string {
	errString := "Max connection attempts exceeded"
	if e.underlyingError != nil {
		errString = fmt.Sprintf("%s - %s", errString, e.underlyingError.Error())
	}
	return errString
}

// WhatsappSender exposes the common operations required for sending whatsapp messages.
// type WhatsappSender interface {
// 	SendWhatsappMessage(from string, to []string, msg io.WriterTo) error
// 	Close() error
// 	Reset() error
// }

// Message is an interface that handles the common operations for whatsapp messages
type Message interface {
	Backoff(reason error) error
	Error(err error) error
	Success() error
	GetAuthToken() (string, error)
	GetNumberId() (string, error)
	GenerateMessage() ([]byte, error)
	GetDestNumber() (string, error)
}

// MessageWorker is the worker that receives slices of messages
// on a channel to send. It's assumed that every slice of messages received is meant
// to be sent to the same server.
type MessageWorker struct {
	queue chan []Message
}

// NewMessageWorker returns an instance of MessageWorker with the message queue
// initialized.
func NewMessageWorker() *MessageWorker {
	return &MessageWorker{
		queue: make(chan []Message),
	}
}

// StartMessaging launches the message worker to begin listening on the MessageQueue channel
// for new slices of Messages instances to process.
func (mw *MessageWorker) StartMessaging(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ms := <-mw.queue:
			go func(ctx context.Context, ms []Message) {
				sendMessage(ctx, ms)
			}(ctx, ms)
		}
	}
}

// MessageQueue sends the provided message to the internal queue for processing.
func (mw *MessageWorker) MessageQueue(ms []Message) {
	mw.queue <- ms
}

// errorMessage is a helper to handle erroring out a slice of Message instances
// in the case that an unrecoverable error occurs.
func errorMessage(err error, ms []Message) {
	for _, m := range ms {
		m.Error(err)
	}
}

// sendMessage attempts to send the provided Message instances.
// If the context is cancelled before all of the messages are sent,
// sendMessage just returns and does not modify those messages.
func sendMessage(ctx context.Context, ms []Message) {
	httpClient := &http.Client{}
	for _, m := range ms {
		select {
		case <-ctx.Done():
			return
		default:
			break
		}

		number_id, err := m.GetNumberId()
		if err != nil {
			m.Error(err)
			continue
		}

		dest_number, err := m.GetDestNumber()
		if err != nil {
			m.Error(err)
			continue
		}

		WhatsappAPIEndpoint := fmt.Sprintf(WhatsappAPI, number_id)

		RequestBody, err := m.GenerateMessage()
		if err != nil {
			m.Error(err)
			continue
		}

		req, err := http.NewRequest("POST", WhatsappAPIEndpoint, bytes.NewBuffer(RequestBody))
		if err != nil {
			log.WithFields(logrus.Fields{
				"code":   "unknown",
				"number": dest_number,
			}).Warn(err)
			errorMessage(err, ms)
			m.Backoff(err)
			continue
		}

		auth_token, err := m.GetAuthToken()
		if err != nil {
			log.WithFields(logrus.Fields{
				"code":   "unknown",
				"number": dest_number,
			}).Warn(err)
			errorMessage(err, ms)
			m.Backoff(err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth_token))
		resp_obj, err := httpClient.Do(req)

		if err != nil {
			log.WithFields(logrus.Fields{
				"code":   "unknown",
				"number": dest_number,
			}).Warn(err)
			errorMessage(err, ms)
			m.Backoff(err)
			continue
		}

		defer resp_obj.Body.Close()
		resp_str, err := io.ReadAll(resp_obj.Body)

		if err != nil {
			m.Error(err)
			continue
		}

		var resp map[string]interface{}
		err = json.Unmarshal([]byte(resp_str), &resp)

		if err != nil {
			m.Error(err)
			continue
		}

		error_resp, ok := resp["error"]
		var error map[string]interface{}

		if !ok {
			_ = json.Unmarshal([]byte(error_resp.(string)), &error)
			err = ErrWhatsappAPIResponse
			log.WithFields(logrus.Fields{
				"code":   error["code"],
				"number": dest_number,
			}).Warn(err)
			errorMessage(err, ms)
			m.Backoff(err)
			continue
		}

		log.WithFields(logrus.Fields{
			"numer_id":    number_id,
			"dest_number": dest_number,
		}).Info("Message sent")
		m.Success()
	}
}
