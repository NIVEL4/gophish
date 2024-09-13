package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/url"
	"path"
	"strings"
	"text/template"

	log "github.com/gophish/gophish/logger"
	qrcode "github.com/skip2/go-qrcode" //library for generating qrcode
)

// TemplateContext is an interface that allows both campaigns and email
// requests to have a PhishingTemplateContext generated for them.
type TemplateContext interface {
	getFromAddress() string
	getBaseURL() string
}

// PhishingTemplateContext is the context that is sent to any template, such
// as the email or landing page content.
type PhishingTemplateContext struct {
	From        string
	URL         string
	Tracker     string
	TrackingURL string
	RId         string
	BaseURL     string
	QrCodeHTML  string
	QrCodeB64   string
	BaseRecipient
}

// NewPhishingTemplateContext returns a populated PhishingTemplateContext,
// parsing the correct fields from the provided TemplateContext and recipient.
func NewPhishingTemplateContext(ctx TemplateContext, r BaseRecipient, rid string) (PhishingTemplateContext, error) {
	f, err := mail.ParseAddress(ctx.getFromAddress())
	if err != nil {
		return PhishingTemplateContext{}, err
	}
	fn := f.Name
	if fn == "" {
		fn = f.Address
	}
	templateURL, err := ExecuteTemplate(ctx.getBaseURL(), r)
	if err != nil {
		return PhishingTemplateContext{}, err
	}

	// For the base URL, we'll reset the the path and the query
	// This will create a URL in the form of http://example.com
	baseURL, err := url.Parse(templateURL)
	if err != nil {
		return PhishingTemplateContext{}, err
	}
	baseURL.Path = ""
	baseURL.RawQuery = ""

	phishURL, _ := url.Parse(templateURL)
	q := phishURL.Query()
	q.Set(RecipientParameter, rid)
	phishURL.RawQuery = q.Encode()

	trackingURL, _ := url.Parse(templateURL)
	trackingURL.Path = path.Join(trackingURL.Path, "/track")
	trackingURL.RawQuery = q.Encode()

	qr_conf, err := GetQR()
	if err != nil {
		log.Error(err)
	}

	var qr *qrcode.QRCode
	qr, _ = qrcode.New(phishURL.String(), qrcode.Medium)

	qrCodeHtml := generateQRCodeHTML(qr, qr_conf)
	qrCodeB64 := generateQRCodeB64(qr, qr_conf)

	return PhishingTemplateContext{
		BaseRecipient: r,
		BaseURL:       baseURL.String(),
		URL:           phishURL.String(),
		TrackingURL:   trackingURL.String(),
		Tracker:       "<img alt='' style='display: none' src='" + trackingURL.String() + "'/>",
		From:          fn,
		RId:           rid,
		QrCodeHTML:    qrCodeHtml,
		QrCodeB64:     qrCodeB64,
	}, nil
}

// ExecuteTemplate creates a templated string based on the provided
// template body and data.
func ExecuteTemplate(text string, data interface{}) (string, error) {
	buff := bytes.Buffer{}
	tmpl, err := template.New("template").Parse(text)
	if err != nil {
		return buff.String(), err
	}
	err = tmpl.Execute(&buff, data)
	return buff.String(), err
}

// ValidationContext is used for validating templates and pages
type ValidationContext struct {
	FromAddress string
	BaseURL     string
}

func (vc ValidationContext) getFromAddress() string {
	return vc.FromAddress
}

func (vc ValidationContext) getBaseURL() string {
	return vc.BaseURL
}

// ValidateTemplate ensures that the provided text in the page or template
// uses the supported template variables correctly.
func ValidateTemplate(text string) error {
	vc := ValidationContext{
		FromAddress: "foo@bar.com",
		BaseURL:     "http://example.com",
	}
	td := Result{
		BaseRecipient: BaseRecipient{
			Email:     "foo@bar.com",
			FirstName: "Foo",
			LastName:  "Bar",
			Position:  "Test",
		},
		RId: "123456",
	}
	ptx, err := NewPhishingTemplateContext(vc, td.BaseRecipient, td.RId)
	if err != nil {
		return err
	}
	_, err = ExecuteTemplate(text, ptx)
	if err != nil {
		return err
	}
	return nil
}

// Generate Qrcode HTML representation
func generateQRCodeHTML(qr *qrcode.QRCode, qr_conf QR) string {
	qrCode := qr.Bitmap()

	// Determine QR code dimensions
	qrWidth := len(qrCode)

	pixelSize := int(qr_conf.Size) / qrWidth

	// Construct HTML table
	var html strings.Builder
	tableOpen := fmt.Sprintf("<table style=\"border-collapse: collapse; border: none; background-color: %s;\">", qr_conf.Background)
	html.WriteString(tableOpen)

	for y := 0; y < qrWidth; y++ {
		html.WriteString("<tr>")
		x := 0
		for x < qrWidth {
			colspan := 0
			if qrCode[y][x] {
				for x < qrWidth && qrCode[y][x] {
					colspan++
					x++
				}
				qrForeground := fmt.Sprintf("<td style=\"width: %dpx; height: %dpx; background-color: %s; border: 0px;\" colspan=%d></td>", pixelSize, pixelSize, qr_conf.Pixels, colspan)
				html.WriteString(qrForeground)
			} else {
				for x < qrWidth && !qrCode[y][x] {
					colspan++
					x++
				}
				qrBackground := fmt.Sprintf("<td style=\"width: %dpx; height: %dpx; border: none;\" colspan=%d></td>", pixelSize, pixelSize, colspan)
				html.WriteString(qrBackground)
			}
		}
		html.WriteString("</tr>")
	}

	html.WriteString("</table>\n")

	return html.String()
}

func generateQRCodeB64(qr *qrcode.QRCode, qr_conf QR) string {
	qr.BackgroundColor, _ = qr_conf.GetBackgroundColor()
	qr.ForegroundColor, _ = qr_conf.GetForegroundColor()

	log.Info(qr)

	// QR to PNG
	QRPNG, _ := qr.PNG(int(qr_conf.Size))

	// B64 encoding PNG
	QRb64 := base64.StdEncoding.EncodeToString(QRPNG)

	// Construct img element
	var img strings.Builder
	img.WriteString("<img src=\"data:image/png;base64,")
	img.WriteString(QRb64)
	img.WriteString("\">")

	return img.String()
}
