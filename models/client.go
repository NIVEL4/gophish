package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/jinzhu/gorm"
)

// Client contains the client data
type Client struct {
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Monitor_url      string     `json:"monitor_url"`
	Monitor_password string     `json:"monitor_password"`
	Apolo_api_key    string     `json:"apolo_api_key"`
	Created_at       time.Time  `json:"created_at"`
	Sent_date        *time.Time `json:"sent_date"`
	Sent_by          *string    `json:"sent_by"`
	Send_method      *string    `json:"send_method"`
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
	Sent_date        *time.Time `json:"sent_date"`
	Sent_by          *string    `json:"sent_by"`
	Send_method      *string    `json:"send_method"`
}

// TableName specifies the database table name for Gorm
func (client *Client) TableName() string {
	return "client"
}

// TableName specifies the database table name for Gorm
func (clientHistory *ClientHistory) TableName() string {
	return "client_history"
}

// GetClient retrieves the latest client data
func GetClient() (Client, error) {
	var client Client
	err := db.Order("created_at desc").First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Info("No client found. Returning empty client.")
		return Client{}, nil
	}
	if err != nil {
		log.Error("Error retrieving client: ", err)
		return Client{}, err
	}
	return client, nil
}

// GetAllClientHistory retrieves all records from the client_history table (including the current client)
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

// Validate performs validation on the client data
func (client *Client) Validate() error {
	return nil
}

// Save client changes or insert a new Client, also records it in the client history
func SaveClient(client *Client) error {
	client.Created_at = time.Now()

	if client.Sent_by == nil {
		defaultSentBy := "Not sent"
		client.Sent_by = &defaultSentBy
	}
	if client.Send_method == nil {
		defaultSendMethod := "None"
		client.Send_method = &defaultSendMethod
	}

	err := db.Create(client).Error
	if err != nil {
		log.Error("Error inserting client: ", err)
		return err
	}

	clientHistory := ClientHistory{
		Name:             client.Name,
		Email:            client.Email,
		Monitor_url:      client.Monitor_url,
		Monitor_password: client.Monitor_password,
		Apolo_api_key:    client.Apolo_api_key,
		Created_at:       client.Created_at,
		Sent_date:        client.Sent_date,
		Sent_by:          client.Sent_by,
		Send_method:      client.Send_method,
	}

	err = db.Create(&clientHistory).Error
	if err != nil {
		log.Error("Error inserting into client_history: ", err)
	}

	return err
}

// UpdateClient updates the client data and logs the change in client_history
func UpdateClient(client *Client) error {
	var existingClient Client

	err := db.Order("created_at desc").First(&existingClient).Error
	if err != nil {
		log.Info("No existing client found. Creating a new client instead.")
		return SaveClient(client)
	}

	if client.Sent_by == nil {
		defaultSentBy := "Not sent"
		client.Sent_by = &defaultSentBy
	}
	if client.Send_method == nil {
		defaultSendMethod := "None"
		client.Send_method = &defaultSendMethod
	}
	if client.Sent_date == nil {
		now := time.Now()
		client.Sent_date = &now
	}

	// This ensures that sent_date, send_method, and sent_by are cleared if the client's name changes.
	if client.Name != existingClient.Name {
		client.Sent_date = nil
		client.Send_method = nil
		client.Sent_by = nil
	}

	// Assigns the current timestamp to sent_date when sent via "Gophish" or "Apolo".
	if client.Send_method != nil && (*client.Send_method == "Gophish" || *client.Send_method == "Apolo") {
		if client.Sent_date == nil {
			now := time.Now()
			client.Sent_date = &now
		}
	}

	err = db.Model(&existingClient).Updates(client).Error
	if err != nil {
		log.Error("Error updating client: ", err)
		return err
	}

	// Store the updated client state in client_history
	clientHistory := ClientHistory{
		Name:             existingClient.Name,
		Email:            existingClient.Email,
		Monitor_url:      existingClient.Monitor_url,
		Monitor_password: existingClient.Monitor_password,
		Apolo_api_key:    existingClient.Apolo_api_key,
		Created_at:       time.Now(),
		Sent_date:        client.Sent_date,
		Sent_by:          client.Sent_by,
		Send_method:      client.Send_method,
	}

	// Insert the updated client state into `client_history`
	err = db.Create(&clientHistory).Error
	if err != nil {
		log.Error("Error inserting into client_history: ", err)
		return err
	}

	return nil
}
