package mailer

import (
	"fmt"
	"strings"

	"github.com/gophish/gomail"
	log "github.com/gophish/gophish/logger"
)

// MonitorEmailRequest defines the expected structure for the email request.
type MonitorEmailRequest struct {
	ClientName        string `json:"client_name"`
	ClientEmail       string `json:"client_email"`
	ClientMonitorURL  string `json:"client_monitor_url"`
	ClientMonitorPass string `json:"client_monitor_password"`
	ClientAPIKey      string `json:"client_api_key"`
	SpecialistName    string `json:"specialist_name"`
	EmailTemplate     string `json:"email_template"`
}

// Send Phishing Monitor Email sends an email to client using the SMTP dialer via gophish.
// SendPhishingMonitorEmail sends an email to the client using the SMTP dialer via Gophish.
func SendPhishingMonitorEmail(dialer Dialer, emailData MonitorEmailRequest, fromAddress string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromAddress)
	m.SetHeader("To", emailData.ClientEmail)
	m.SetHeader("Subject", "Phishing Monitor Details")

	// Format the URL to replace dots with [.]
	formattedURL := strings.NewReplacer(
		".", "[.]",
		":", "[:]",
	).Replace(emailData.ClientMonitorURL)

	// Define email body based on the selected template
	var body string
	switch emailData.EmailTemplate {
	case "1":
		// Option 1: Only Monitor Password
		body = fmt.Sprintf(`
		<html>
		<body>
			<p>Dear %s,</p>
			<p>Here are your monitoring access:</p>
			<p><b>Access Password:</b> %s</p>
			<p>Best regards,</p>
			<p>%s.</p>
		</body>
		</html>
		`, emailData.ClientName, emailData.ClientMonitorPass, emailData.SpecialistName)
	case "2":
		// Option 2: Panel and Password
		body = fmt.Sprintf(`
		<html>
		<body>
			<p>Dear %s,</p>
			<p>We are providing you with access to the panel.</p>
			<p><b>Monitoring Panel:</b> %s</p>
			<p><b>Access Password:</b> %s</p>
			<p>Best regards,</p>
			<p>%s.</p>
		</body>
		</html>
		`, emailData.ClientName, formattedURL, emailData.ClientMonitorPass, emailData.SpecialistName)
	default:
		// Default case in case of an invalid option
		body = "<p>Invalid email template selected.</p>"
	}

	m.SetBody("text/html", body)

	// Connect to the SMTP server using the dialer
	sender, err := dialer.Dial()
	if err != nil {
		log.Error("Error connecting to the SMTP server:", err)
		return err
	}

	// Send the email
	if err := gomail.Send(sender, m); err != nil {
		log.Error("Error sending email:", err)
		return err
	}

	log.Info("Email successfully sent to:", emailData.ClientEmail)
	return nil
}
