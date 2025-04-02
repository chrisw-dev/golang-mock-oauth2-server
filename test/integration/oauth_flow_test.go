package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/server"
)

func TestOAuthFlow(t *testing.T) {
	// Create and start the server
	mockServer := server.NewServer(":0") // Use port 0 to get a random available port

	// Start server in a goroutine
	ts := httptest.NewServer(mockServer.Handler())
	defer ts.Close()

	// Create OAuth2 config pointing to test server
	oauth2Config := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/authorize",
			TokenURL: ts.URL + "/token",
		},
	}

	// Test the authorization URL
	authURL := oauth2Config.AuthCodeURL("test-state")
	resp, err := http.Get(authURL)
	if err != nil {
		t.Fatalf("Failed to get auth URL: %v", err)
	}
	defer resp.Body.Close()

	// Check redirect to callback with code
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected redirect, got status: %d", resp.StatusCode)
	}

	// Extract code from redirect URL
	location := resp.Header.Get("Location")
	if location == "" {
		t.Fatalf("No Location header in response")
	}
	// Parse the code from the URL...

	// Test exchanging the code for a token
	// Test accessing the user info endpoint
}
