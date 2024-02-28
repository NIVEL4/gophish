package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gophish/gophish/config"
	"github.com/gophish/gophish/models"
)

type testContext struct {
	apiKey    string
	config    *config.Config
	apiServer *Server
	admin     models.User
}

func setupTest(t *testing.T) *testContext {
	conf := &config.Config{
		DBName:         "sqlite3",
		DBPath:         ":memory:",
		MigrationsPath: "../../db/db_sqlite3/migrations/",
	}
	err := models.Setup(conf)
	if err != nil {
		t.Fatalf("Failed creating database: %v", err)
	}
	ctx := &testContext{}
	ctx.config = conf
	// Get the API key to use for these tests
	u, err := models.GetUser(1)
	if err != nil {
		t.Fatalf("error getting admin user: %v", err)
	}
	ctx.apiKey = u.ApiKey
	ctx.admin = u
	ctx.apiServer = NewServer()
	return ctx
}

func createTestData(t *testing.T) (string, *testContext) {
	ctx := setupTest(t) // Set up the test context

	// Add a group
	group := models.Group{Name: "Test Group"}
	group.Targets = []models.Target{
		{BaseRecipient: models.BaseRecipient{Email: "test1@example.com", FirstName: "First", LastName: "Example"}},
		{BaseRecipient: models.BaseRecipient{Email: "test2@example.com", FirstName: "Second", LastName: "Example"}},
	}
	group.UserId = 1
	if err := models.PostGroup(&group); err != nil {
		t.Fatalf("error adding group: %v", err)
	}

	// Add a template
	template := models.Template{Name: "Test Template"}
	template.Subject = "Test subject"
	template.Text = "Text text"
	template.HTML = "<html>Test</html>"
	template.UserId = 1
	if err := models.PostTemplate(&template); err != nil {
		t.Fatalf("error adding template: %v", err)
	}

	// Add a landing page
	p := models.Page{Name: "Test Page"}
	p.HTML = "<html>Test</html>"
	p.UserId = 1
	if err := models.PostPage(&p); err != nil {
		t.Fatalf("error adding landing page: %v", err)
	}

	// Add a sending profile
	smtp := models.SMTP{Name: "Test Page"}
	smtp.UserId = 1
	smtp.Host = "example.com"
	smtp.FromAddress = "test@test.com"
	if err := models.PostSMTP(&smtp); err != nil {
		t.Fatalf("error adding sending profile: %v", err)
	}

	// Setup and "launch" our campaign
	// Set the status such that no emails are attempted
	c := models.Campaign{Name: "Test campaign"}
	c.UserId = 1
	c.Template = template
	c.Page = p
	c.SMTP = smtp
	c.Groups = []models.Group{group}
	if err := models.PostCampaign(&c, c.UserId); err != nil {
		t.Fatalf("error setting up campaign: %v", err)
	}
	if err := c.UpdateStatus(models.CampaignEmailsSent); err != nil {
		t.Fatalf("error updating campaign status: %v", err)
	}

	// Return the expected HTML response and the test context
	return "<html><head><base href=\"\"/></head><body><img src=\"/test.png\"/></body></html>", ctx
}

func TestSiteImportBaseHref(t *testing.T) {
	expected, ctx := createTestData(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expected)
	}))
	defer ts.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/import/site",
		bytes.NewBuffer([]byte(fmt.Sprintf(`
            {
                "url" : "%s",
                "include_resources" : false
            }
        `, ts.URL))))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ctx.apiServer.ImportSite(response, req)

	var cs cloneResponse
	err := json.NewDecoder(response.Body).Decode(&cs)
	if err != nil {
		t.Fatalf("error decoding response: %v", err)
	}
	fmt.Println("Actual:")
	fmt.Println(cs.HTML)
	fmt.Println("Expected:")
	fmt.Println(expected)
	cs.HTML = strings.ReplaceAll(cs.HTML, "\n", "")

	if cs.HTML != expected {
		t.Fatalf("unexpected response received. expected %s got %s", expected, cs.HTML)
	}

}
