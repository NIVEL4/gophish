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
	m.SetHeader("Subject", "Acceso al Panel de Monitoreo de Phishing")

	// Format the URL to replace dots with [.]
	formattedURL := strings.NewReplacer(
		".", "[.]",
		":", "[:]",
	).Replace(emailData.ClientMonitorURL)

	// Define email body based on the selected template
	var body string
	switch emailData.EmailTemplate {
	case "1":
		// Option 1: Panel and Password
		body = fmt.Sprintf(`
		<html>
		<body>
			<p>Estimado/a %s,</p>
			<p>Le enviamos sus credenciales de acceso para el panel de monitoreo:</p>
			<p><b>Panel de monitoreo:</b> %s</p>
			<p><b>Contraseña de acceso:</b> %s</p>
			<p>Le agradeceríamos que confirme su acceso al panel una vez haya logrado ingresar correctamente, dejando constancia en el hilo principal.</p>
			<p><em>Este es un mensaje automático y no debe responderlo.</em></p>
			<p>Atentamente,</p>  
			<p><strong>%s</strong></p>
		</body>
		</html>
		`, emailData.ClientName, formattedURL, emailData.ClientMonitorPass, emailData.SpecialistName)
	case "2":
		// Option 2: Only Monitor Password
		body = fmt.Sprintf(`
		<html>
		<body>
		<p>Estimado/a %s,</p>
		<p>Le enviamos sus credenciales de acceso para el panel de monitoreo:</p>
		<p><strong>Contraseña de acceso:</strong> %s</p>
		<p>Le agradeceríamos que confirme su acceso al panel una vez haya logrado ingresar correctamente, dejando constancia en el hilo principal.</p>
		<p><em>Este es un mensaje automático y no debe responderlo.</em></p>
		<p>Atentamente,</p>
		<p><strong>%s</strong></p>
		</body>
		</html>
		`, emailData.ClientName, emailData.ClientMonitorPass, emailData.SpecialistName)
	default:
		// Default case in case of an invalid option
		body = "<p>Por favor, solicite el correo nuevamente al especialista.</p>"
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
