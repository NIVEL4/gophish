package models

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net"
	"sync"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/jinzhu/gorm"
	"github.com/oschwald/maxminddb-golang"
)

type mmCity struct {
	GeoPoint mmGeoPoint `maxminddb:"location"`
}

type mmGeoPoint struct {
	Latitude  float64 `maxminddb:"latitude"`
	Longitude float64 `maxminddb:"longitude"`
}

// Result contains the fields for a result object,
// which is a representation of a target in a campaign.
type Result struct {
	Id           int64     `json:"-"`
	CampaignId   int64     `json:"-"`
	UserId       int64     `json:"-"`
	RId          string    `json:"id"`
	Status       string    `json:"status" sql:"not null"`
	IP           string    `json:"ip"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	SendDate     time.Time `json:"send_date"`
	Reported     bool      `json:"reported" sql:"not null"`
	ModifiedDate time.Time `json:"modified_date"`
	BaseRecipient
}

// ResultSemaphor allows to manage access to critical DB operations
// avoiding race conditions
type ResultSemaphor struct {
	locked 	map[int]bool{}
	mu sync.Mutex
}

var resultSemaphor ResultSemaphor

// waitForTurn waits until no other goroutine is accessing 
// a specific result (Id = resultId) to continue
func waitForTurn(resultId) {
	var myTurn := false
	for !myTurn {
		resultSemaphor.mu.Lock()
		_, ok := resultSemaphor.locked[resultId]
		if !ok {
			resultSemaphor.locked[resultId] = true
			myTurn = true
		}
		resultSemaphor.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

// releaseTurn releases the turn for a specific result
// (Id = resultId), so other goroutines can access it
func releaseTurn(resultId) {
	resultSemaphor.mu.Lock()
	delete(resultSemaphor.locked, resultId)
	resultSemaphor.mu.Unlock()
}

func (r *Result) createEvent(status string, details interface{}) (*Event, error) {
	e := &Event{Email: r.Email, Message: status}
	if details != nil {
		dj, err := json.Marshal(details)
		if err != nil {
			return nil, err
		}
		e.Details = string(dj)
	}
	AddEvent(e, r.CampaignId)
	return e, nil
}

// HandleEmailSent updates a Result to indicate that the email has been
// successfully sent to the remote SMTP server
func (r *Result) HandleEmailSent() error {
	event, err := r.createEvent(EventSent, nil)
	if err != nil {
		return err
	}
	waitForTurn(r.Id)
	r.SendDate = event.Time
	r.Status = EventSent
	r.ModifiedDate = event.Time
	err = db.Save(r).Error
	releaseTurn(r.Id)
	return err
}

// HandleEmailError updates a Result to indicate that there was an error when
// attempting to send the email to the remote SMTP server.
func (r *Result) HandleEmailError(err error) error {
	event, err := r.createEvent(EventSendingError, EventError{Error: err.Error()})
	if err != nil {
		return err
	}
	r.Status = Error
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// HandleEmailBackoff updates a Result to indicate that the email received a
// temporary error and needs to be retried
func (r *Result) HandleEmailBackoff(err error, sendDate time.Time) error {
	event, err := r.createEvent(EventSendingError, EventError{Error: err.Error()})
	if err != nil {
		return err
	}
	r.Status = StatusRetry
	r.SendDate = sendDate
	r.ModifiedDate = event.Time
	return db.Save(r).Error
}

// HandleEmailOpened updates a Result in the case where the recipient opened the
// email.
func (r *Result) HandleEmailOpened(details EventDetails) error {
	event, err := r.createEvent(EventOpened, details)
	if err != nil {
		return err
	}
	// Don't update the status if the user already opened the email, 
	// clicked the link or submitted data to the campaign
	waitForTurn(r.Id)
	if r.Status == EventClicked || r.Status == EventDataSubmit || r.Status == EventOpened {
		err = nil
	} else {
		r.Status = EventOpened
		r.ModifiedDate = event.Time
		err = db.Save(r).Error
	}
	releaseTurn(r.Id)
	return err
}

// HandleClickedLink updates a Result in the case where the recipient clicked
// the link in an email.
func (r *Result) HandleClickedLink(details EventDetails) error {
	event, err := r.createEvent(EventClicked, details)
	if err != nil {
		return err
	}
	// Don't update the status if the user has already clicked the link or 
	// submitted data via the landing page form.
	waitForTurn(r.Id)
	if r.Status == EventDataSubmit || r.Status == EventClicked {
		err = nil
	} else {
		r.Status = EventClicked
		r.ModifiedDate = event.Time
		err = db.Save(r).Error
	}
	releaseTurn(r.Id)
	return err
}

// HandleFormSubmit updates a Result in the case where the recipient submitted
// credentials to the form on a Landing Page.
func (r *Result) HandleFormSubmit(details EventDetails) error {
	event, err := r.createEvent(EventDataSubmit, details)
	if err != nil {
		return err
	}
	// Don't update the status if the user has already submitted data
	// voia the landing page form.
	waitForTurn(r.Id)
	if r.Status == EventDataSubmit {
		err = nil
	} else {
		r.Status = EventDataSubmit
		r.ModifiedDate = event.Time
		err = db.Save(r).Error
	}
	releaseTurn(r.Id)
	return err
}

// HandleEmailReport updates a Result in the case where they report a simulated
// phishing email using the HTTP handler.
func (r *Result) HandleEmailReport(details EventDetails) error {
	event, err := r.createEvent(EventReported, details)
	if err != nil {
		return err
	}
	waitForTurn(r.Id)
	r.Reported = true
	r.ModifiedDate = event.Time
	err = db.Save(r).Error
	releaseTurn(r.Id)
	return err
}

// UpdateGeo updates the latitude and longitude of the result in
// the database given an IP address
func (r *Result) UpdateGeo(addr string) error {
	// Open a connection to the maxmind db
	mmdb, err := maxminddb.Open("static/db/geolite2-city.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer mmdb.Close()
	ip := net.ParseIP(addr)
	var city mmCity
	// Get the record
	err = mmdb.Lookup(ip, &city)
	if err != nil {
		return err
	}
	// Update the database with the record information
	waitForTurn(r.Id)
	r.IP = addr
	r.Latitude = city.GeoPoint.Latitude
	r.Longitude = city.GeoPoint.Longitude
	err = db.Save(r).Error
	releaseTurn(r.Id)
	return err
}

func generateResultId() (string, error) {
	const alphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	k := make([]byte, 7)
	for i := range k {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNum))))
		if err != nil {
			return "", err
		}
		k[i] = alphaNum[idx.Int64()]
	}
	return string(k), nil
}

// GenerateId generates a unique key to represent the result
// in the database
func (r *Result) GenerateId(tx *gorm.DB) error {
	// Keep trying until we generate a unique key (shouldn't take more than one or two iterations)
	for {
		rid, err := generateResultId()
		if err != nil {
			return err
		}
		r.RId = rid
		err = tx.Table("results").Where("r_id=?", r.RId).First(&Result{}).Error
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	return nil
}

// GetResult returns the Result object from the database
// given the ResultId
func GetResult(rid string) (Result, error) {
	r := Result{}
	err := db.Where("r_id=?", rid).First(&r).Error
	return r, err
}
