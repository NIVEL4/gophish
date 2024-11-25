package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/gophish/gomail"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/mailer"
)

// Dialer is a wrapper around a standard gomail.Dialer in order
// to implement the mailer.Dialer interface. This allows us to better
// separate the mailer package as opposed to forcing a connection
// between mailer and gomail.
type WhatsappDialer struct {
	*gomail.Dialer
}

// Dial wraps the gomail dialer's Dial command
func (d *Dialer) DialWhatsapp() (mailer.Sender, error) {
	return d.Dialer.Dial()
}

// SMTP contains the attributes needed to handle the sending of campaign emails
type Whatsapp struct {
	Id           int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64     `json:"-" gorm:"column:user_id"`
	Interface    string    `json:"interface_type" gorm:"column:interface_type"`
	Name         string    `json:"name"`
	Number       string    `json:"number"`
	AuthToken    string    `json:"auth_token"`
	ModifiedDate time.Time `json:"modified_date"`
}

// ErrNumberNotSpecified is thrown when there is no "Number"
// specified in the Whatsapp configuration
var ErrNumberNotSpecified = errors.New("No Number specified")

// ErrInvalidNumber is thrown when the Whatsapp Number field in the sending
// profiles containes a value that is not a Whatsapp number
var ErrInvalidNumber = errors.New("Invalid Number because it is not a whatsapp number")

// ErrAuthTokenNotSpecified is thrown when there is no Auth Token specified
// in the Whatsapp configuration
var ErrAuthTokenNotSpecified = errors.New("No Auth Token specified")

// ErrInvalidAuthToken indicates that the Auth Token string is invalid
var ErrInvalidAuthToken = errors.New("Invalid Auth Token")

// TableName specifies the database tablename for Gorm to use
func (w Whatsapp) TableName() string {
	return "whatsapp"
}

// Validate ensures that SMTP configs/connections are valid
func (w *Whatsapp) Validate() error {
	switch {
	case w.Number == "":
		return ErrNumberNotSpecified
	case w.AuthToken == "":
		return ErrAuthTokenNotSpecified
	case !ValidateWhatsappNumber(w.Number):
		return ErrInvalidNumber
	}
	return nil
}

// validateFromAddress validates
func ValidateWhatsappNumber(number string) bool {
	r, _ := regexp.Compile("^[+][0-9]{2}[ ]?([0-9]+(-| )?)+$")
	return r.MatchString(number)
}

// GetWhatsapps returns the list of Whatsapp owned by the given user.
func GetWhatsapps(uid int64) ([]Whatsapp, error) {
	ws := []Whatsapp{}
	err := db.Where("user_id=?", uid).Find(&ws).Error
	if err != nil {
		log.Error(err)
		return ws, err
	}
	return ws, nil
}

// GetWhatsapp returns the Whatsapp, if it exists, specified by the given id and user_id.
func GetWhatsapp(id int64, uid int64) (Whatsapp, error) {
	w := Whatsapp{}
	err := db.Where("user_id=? and id=?", uid, id).Find(&w).Error
	if err != nil {
		log.Error(err)
		return w, err
	}
	return w, err
}

// GetWhatsappByName returns the Whatsapp, if it exists, specified by the given name and user_id.
func GetWhatsappByName(n string, uid int64) (Whatsapp, error) {
	w := Whatsapp{}
	err := db.Where("user_id=? and name=?", uid, n).Find(&w).Error
	if err != nil {
		log.Error(err)
		return w, err
	}
	return w, err
}

// PostWhatsapp creates a new Whatsapp in the database.
func PostWhatsapp(w *Whatsapp) error {
	err := w.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	// Insert into the DB
	err = db.Save(w).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// PutWhatsapp edits an existing Whatsapp in the database.
// Per the PUT Method RFC, it presumes all data for a Whatsapp is provided.
func PutWhatsapp(w *Whatsapp) error {
	err := w.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	err = db.Where("id=?", w.Id).Save(w).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// DeleteWhatsapp deletes an existing Whatsapp in the database.
// An error is returned if a Whatsapp with the given user id and
// Whatsapp id is not found.
func DeleteWhatsapp(id int64, uid int64) error {
	err := db.Where("user_id=?", uid).Delete(Whatsapp{Id: id}).Error
	if err != nil {
		log.Error(err)
	}
	return err
}
