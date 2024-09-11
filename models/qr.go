package models

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// QR contains the settings for generating QR codes
type QR struct {
	UserId		 int64	   `json:"user_id"`
	Size         int64     `json:"size"`
	Pixels       string    `json:"pixels"`
	Background   string    `json:"background"`
}

// ErrQRCodeTooSmall is thrown when QR code dimensions are too small for a QR code
var ErrQRCodeTooSmall = errors.New("QR Code is too small")

// ErrInvalidColor is thrown when an invalid color is chosen for a QR code
var ErrInvalidColor = errors.New("Invalid color")

// Validate performs validation on given settings
func (qr *QR) Validate() error {
	if qr.Size < 64 {
		return ErrQRCodeTooSmall
	}
	r, _ := regexp.MustCompile("^#[[:xdigit:]]{3,6}$")
	if !r.MatchString(qr.Pixels) || !r.MatchString(qr.Background) {
		return ErrInvalidColor
	}
	return nil
}

// TableName specifies the database tablename for Gorm to use
func (qr QR) TableName() string {
	return "qr_conf"
}

// GetQR returns the QR settings
func GetQR(uid int64) (QR, error) {
	qr := QR{}
	err := db.Where("user_id=?", uid).Find(&qr).Error
	if err != nil {
		log.Error(err)
	}
	return g, err
}

// DeleteQR deletes que QR code settings
func DeleteQR(uid int64) error {
	qr := QR{}
	err := db.Where("user_id=?", uid).Find(&qr).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// UpdateQR updates the QR code settings
func UpdateQR(qr *QR, uid int64) error {
	err := qr.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	err = DeleteQR(uid)
	if err != nil {
		log.Error(err)
		return err
	}
	err = db.Save(qr).Error
	if err != nil {
		log.Error(err)
	}
	return err
}
