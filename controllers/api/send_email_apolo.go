package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Struct to match emailData format (settings.js)
type EmailData struct {
	ClientName        string `json:"client_name"`
	ClientEmail       string `json:"client_email"`
	ClientMonitorURL  string `json:"client_monitor_url"`
	ClientMonitorPass string `json:"client_monitor_password"`
	ClientAPIKey      string `json:"client_api_key"`
	SpecialistName    string `json:"specialist_name"`
	SendDate          string `json:"send_date"`
}

// SendMailApoloAPI handles requests for sending email data via apolo.
func (as *Server) SendMailApoloAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data EmailData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Convert struct to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		// URL to send POST request
		postURL := ""

		// Make POST request
		resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error sending request:", err)
			http.Error(w, "Failed to send request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Respond with the status of the external request
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte("Request sent successfully"))

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
