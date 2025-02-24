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
}

// Send Phishing Monitor Email sends an email to client using the SMTP dialer via gophish.
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

	// Mail Template for testing purpose.
	body := fmt.Sprintf(`
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
