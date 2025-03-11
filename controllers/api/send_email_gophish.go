package api

import (
	"encoding/json"
	"net/http"

	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/mailer"
	"github.com/gophish/gophish/models"
)

// SendMailGophish handles the API request to send an email using Gophish.
func (as *Server) SendMailGophish(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define the expected email data structure
	var emailData struct {
		ClientName        string `json:"client_name"`
		ClientEmail       string `json:"client_email"`
		ClientMonitorURL  string `json:"client_monitor_url"`
		ClientMonitorPass string `json:"client_monitor_password"`
		ClientAPIKey      string `json:"client_api_key"`
		SpecialistName    string `json:"specialist_name"`
		SMTPProfileID     int    `json:"smtp_profile"`
		EmailTemplate     string `json:"email_template"`
	}

	// Decode JSON request body into emailData struct
	if err := json.NewDecoder(r.Body).Decode(&emailData); err != nil {
		log.Error("Error decoding JSON:", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate that SMTP Profile ID is provided
	if emailData.SMTPProfileID == 0 {
		log.Error("SMTP Profile ID not provided")
		http.Error(w, "SMTP Profile is required", http.StatusBadRequest)
		return
	}

	// Retrieve SMTP configuration based on the selected profile
	smtpConfig, err := models.GetSMTP(int64(emailData.SMTPProfileID), 1)
	if err != nil {
		log.Error("Error fetching SMTP configuration:", err)
		http.Error(w, "Error retrieving SMTP configuration", http.StatusInternalServerError)
		return
	}

	// Get the SMTP dialer
	dialer, err := smtpConfig.GetDialer()
	if err != nil {
		log.Error("Error obtaining SMTP dialer:", err)
		http.Error(w, "Error retrieving SMTP dialer", http.StatusInternalServerError)
		return
	}

	// Validate required email data before proceeding
	if emailData.ClientEmail == "" || emailData.ClientMonitorURL == "" {
		log.Error("Missing required email data")
		http.Error(w, "Missing required email data", http.StatusBadRequest)
		return
	}

	// Construct the email request object
	monitorEmail := mailer.MonitorEmailRequest{
		ClientName:        emailData.ClientName,
		ClientEmail:       emailData.ClientEmail,
		ClientMonitorURL:  emailData.ClientMonitorURL,
		ClientMonitorPass: emailData.ClientMonitorPass,
		ClientAPIKey:      emailData.ClientAPIKey,
		SpecialistName:    emailData.SpecialistName,
		EmailTemplate:     emailData.EmailTemplate,
	}

	// Send the email only if all validations are successful
	if err := mailer.SendPhishingMonitorEmail(dialer, monitorEmail, smtpConfig.FromAddress); err != nil {
		log.Error("Error sending email:", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Send success response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}
