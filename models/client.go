package models

import (
	"time"

	log "github.com/gophish/gophish/logger"
)

// Client contains the client data
type Client struct {
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Monitor_url      string     `json:"monitor_url"`
	Monitor_password string     `json:"monitor_password"`
	Apolo_api_key    string     `json:"apolo_api_key"`
	Created_at       time.Time  `json:"created_at"`
	Send_date        *time.Time `json:"send_date"`
}

// ClientHistory stores historical changes to the client data
type ClientHistory struct {
	ID               uint       `json:"id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Monitor_url      string     `json:"monitor_url"`
	Monitor_password string     `json:"monitor_password"`
	Apolo_api_key    string     `json:"apolo_api_key"`
	Created_at       time.Time  `json:"created_at"`
	Send_date        *time.Time `json:"send_date"`
	Change_date      time.Time  `json:"change_date"`
}

// TableName specifies the database table name for Gorm
func (client *Client) TableName() string {
	return "client"
}

// TableName specifies the database table name for Gorm
func (clientHistory *ClientHistory) TableName() string {
	return "client_history"
}

// Validate performs validation on the client data
func (client *Client) Validate() error {
	// You can add any necessary validation here
	return nil
}

// GetClient retrieves the latest client data (most recent)
func GetClient() (Client, error) {
	var client Client
	err := db.Order("created_at desc").Limit(1).Find(&client).Error
	if err != nil {
		log.Error("Error retrieving client: ", err)
	}
	return client, err
}

// GetAllClientHistory retrieves all records from the client_history table
func GetAllClientHistory() ([]ClientHistory, error) {
	var history []ClientHistory
	err := db.Find(&history).Error
	if err != nil {
		log.Error("Error retrieving client history: ", err)
	}
	return history, err
}

// DeleteClient deletes the current client data
func DeleteClient() error {
	client := Client{}
	err := db.Find(&client).Error
	if err != nil {
		log.Error("Error finding client to delete: ", err)
		return err
	}
	err = db.Delete(&client).Error
	if err != nil {
		log.Error("Error deleting client: ", err)
	}
	return err
}

// UpdateClient updates the client data and logs the change in client_history
func UpdateClient(client *Client) error {
	// Check if a client exists before updating
	var existingClient Client
	err := db.First(&existingClient).Error
	if err != nil {
		// No existing client found, so just save the new one
		log.Info("No existing client found, saving new client.")
		return db.Create(client).Error
	}

	// Log the current client data in client_history before updating
	clientHistory := ClientHistory{
		Name:             existingClient.Name,
		Email:            existingClient.Email,
		Monitor_url:      existingClient.Monitor_url,
		Monitor_password: existingClient.Monitor_password,
		Apolo_api_key:    existingClient.Apolo_api_key,
		Created_at:       existingClient.Created_at, // Preserve the original creation date
		Send_date:        existingClient.Send_date,
		Change_date:      time.Now(), // Timestamp for when the update happens
	}

	// Insert the current client data into client_history
	err = db.Create(&clientHistory).Error
	if err != nil {
		log.Error("Error inserting into client_history: ", err)
		return err
	}

	// Delete the old client before updating
	err = DeleteClient()
	if err != nil {
		log.Error("Error deleting previous client before update: ", err)
		return err
	}

	// Save the new client data
	err = db.Create(client).Error // Change from Save() to Create() to avoid issues with primary keys
	if err != nil {
		log.Error("Error updating client: ", err)
	}
	return err
}
