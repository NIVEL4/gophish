package models

import (
	"errors"
	"image/color"
	"regexp"
	"strconv"

	log "github.com/gophish/gophish/logger"
)

// QR contains the settings for generating QR codes
type QR struct {
	Size       int64  `json:"qr_size"`
	Pixels     string `json:"qr_pixels"`
	Background string `json:"qr_background"`
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
	r := regexp.MustCompile("^#[[:xdigit:]]{6}$")
	if !r.MatchString(qr.Pixels) || !r.MatchString(qr.Background) {
		return ErrInvalidColor
	}
	return nil
}

// TableName specifies the database tablename for Gorm to use
func (qr *QR) TableName() string {
	return "qr_conf"
}

// GetQR returns the QR settings
func GetQR() (QR, error) {
	var qrcode QR
	err := db.Find(&qrcode).Error
	if err != nil {
		log.Error(err)
	}
	return qrcode, err
}

// GetForegroundColor returns the color for the foreground pixels
func (qr *QR) GetForegroundColor() (color.Color, error) {
	var c color.RGBA
	r, g, b, err := qr.Str2RGB(qr.Pixels)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	c.R = r
	c.G = g
	c.B = b
	c.A = 255

	return c, nil
}

// GetBackgroundColor returns the color for the background pixels
func (qr *QR) GetBackgroundColor() (color.Color, error) {
	var c color.RGBA
	r, g, b, err := qr.Str2RGB(qr.Background)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	c.R = r
	c.G = g
	c.B = b
	c.A = 255

	return c, nil
}

// DeleteQR deletes que QR code settings
func DeleteQR() error {
	qr := QR{}
	err := db.Find(&qr).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// UpdateQR updates the QR code settings
func UpdateQR(qr *QR) error {
	err := qr.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	err = DeleteQR()
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

func (qr *QR) Str2RGB(cstr string) (uint8, uint8, uint8, error) {
	r, err := strconv.ParseUint(cstr[1:3], 16, 32)
	if err != nil {
		log.Error(err)
		return 0, 0, 0, err
	}
	g, err := strconv.ParseUint(cstr[3:5], 16, 32)
	if err != nil {
		log.Error(err)
		return 0, 0, 0, err
	}
	b, err := strconv.ParseUint(cstr[5:7], 16, 32)
	if err != nil {
		log.Error(err)
		return 0, 0, 0, err
	}
	return uint8(r), uint8(g), uint8(b), nil
}
