package models

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/url"
	"path"
	"strings"
	"text/template"

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
	QrCode      string
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

	q, err := qrcode.New(websiteURL, qrcode.Medium)

	if err == nil {
		qrCodeHtml := generateQRCodeHTML(q)
		qrCodeB64 := generateQRCodeB64(q)
	} else {
		qrCodeHtml := nil
		qrCodeB64 := nil
	}

	return PhishingTemplateContext{
		BaseRecipient:	r,
		BaseURL:	baseURL.String(),
		URL:		phishURL.String(),
		TrackingURL:	trackingURL.String(),
		Tracker:	"<img alt='' style='display: none' src='" + trackingURL.String() + "'/>",
		From:		fn,
		RId:		rid,
		QrCodeHTML:	qrCodeHtml,
		QrCodeB64:	qrCodeB64,
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
func generateQRCodeHTML(qrCode []byte) string {
	// Determine QR code dimensions
	qrWidth := len(qrCode)

	// Construct HTML table
	var html strings.Builder
	html.WriteString("<table style=\"border-collapse: collapse; background-color: transparent;\">")

	for y := 0; y < qrWidth; y++ {
		html.WriteString("<tr>")
		for x := 0; x < qrWidth; x++ {
			if qrCode[y][x] {
				html.WriteString("<td style=\"width: 1px; height: 1px; background-color: black; border: 1px solid black;\"></td>")
			} else {
				html.WriteString("<td style=\"width: 1px; height: 1px; border: none;\"></td>")
			}
		}
		html.WriteString("</tr>")
	}

	html.WriteString("</table>\n")

	return html.String()
}

func generateQRCodeB64(qrCode []byte) string {
	// Base64 encoding QR
	QRb64 := base64.StdEncoding.EncodeToString(b)

	// Construct img element
	var img strings.Builder
	img.WriteString("<img src=\"data:image/png;base64,")
	img.WriteString(QRb64)
	img.WriteString("\">")

	return img.String()
}
