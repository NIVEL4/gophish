package models

import (
	log "github.com/gophish/gophish/logger"
)

// Client contains the client data
type Client struct {
	Name string `json:"name"`
}

// TableName specifies the database tablename for Gorm to use
func (client *Client) TableName() string {
	return "client"
}

// Validate performs validation on given settings
func (client *Client) Validate() error {
        return nil
}

// GetClient returns the Client data
func GetClient() (Client, error) {
	var client Client
	err := db.Find(&client).Error
	if err != nil {
		log.Error(err)
	}
	return client, err
}

// DeleteClient deletes the client data
func DeleteClient() error {
	client := Client{}
	err := db.Find(&client).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// UpdateClient updates the client data
func UpdateClient(client *Client) error {
	err := DeleteClient()
	if err != nil {
		log.Error(err)
	}
	err = db.Save(client).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

