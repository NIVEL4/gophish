package models

import (
	"errors"
	log "github.com/gophish/gophish/logger"
	"strings"
)

// BeEF contains the beef data
type BeEF struct {
	URL string `json:"url"`
}

// ErrNotJS indicates the URL does not point to a JS file
var ErrNotJS = errors.New("Hook file is not JS")

// ErrMissingSchema indicates the URL is missing the http[s] schema
var ErrMissingSchema = errors.New("Missing HTTP[s] schema")

// TableName specifies the database tablename for Gorm to use
func (beef *BeEF) TableName() string {
	return "beef"
}

// Validate performs validation on given settings
func (beef *BeEF) Validate() error {
	if !strings.HasSuffix(beef.URL, ".js") {
		return ErrNotJS
	}
	if !strings.HasPrefix(beef.URL, "http://") && !strings.HasPrefix(beef.URL, "https://") {
		return ErrMissingSchema
	}
	return nil
}

// GetBeEF returns the beef data
func GetBeEF() (BeEF, error) {
	var beef BeEF
	err := db.Find(&beef).Error
	if err != nil {
		log.Error(err)
	}
	return beef, err
}

// DeleteBeEF deletes the beef data
func DeleteBeEF() error {
	beef := BeEF{}
	err := db.Find(&beef).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// UpdateBeEF updates the beef data
func UpdateBeEF(beef *BeEF) error {
	err := DeleteClient()
	if err != nil {
		log.Error(err)
	}
	err = db.Save(beef).Error
	if err != nil {
		log.Error(err)
	}
	return err
}
