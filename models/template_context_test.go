package models

import (
	"fmt"
	"os"

	check "gopkg.in/check.v1"
)

type mockTemplateContext struct {
	URL         string
	FromAddress string
}

func (m mockTemplateContext) getFromAddress() string {
	return m.FromAddress
}

func (m mockTemplateContext) getBaseURL() string {
	return m.URL
}

func (s *ModelsSuite) TestNewTemplateContext(c *check.C) {
	r := Result{
		BaseRecipient: BaseRecipient{
			FirstName: "Foo",
			LastName:  "Bar",
			Email:     "foo@bar.com",
		},
		RId: "1234567",
	}
	ctx := mockTemplateContext{
		URL:         "https://example.com",
		FromAddress: "From Address <from@example.com>",
	}

	expectedHTMLFilePath := "./qr-test-data/test-qr-code.html"
	expectedHTMLBytes, err := os.ReadFile(expectedHTMLFilePath)
	if err != nil {
		c.Fatalf("Failed to read new-html.html: %v", err)
	}
	// Remove extra backslashes from the read string
	expectedHTMLString := string(expectedHTMLBytes)

	expected := PhishingTemplateContext{
		URL:           fmt.Sprintf("%s?rid=%s", ctx.URL, r.RId),
		BaseURL:       ctx.URL,
		BaseRecipient: r.BaseRecipient,
		TrackingURL:   fmt.Sprintf("%s/track?rid=%s", ctx.URL, r.RId),
		From:          "From Address",
		RId:           r.RId,
		QrCode:        expectedHTMLString,
	}
	expected.Tracker = "<img alt='' style='display: none' src='" + expected.TrackingURL + "'/>"
	got, err := NewPhishingTemplateContext(ctx, r.BaseRecipient, r.RId)
	c.Assert(err, check.Equals, nil)
	c.Assert(got, check.DeepEquals, expected)
}
