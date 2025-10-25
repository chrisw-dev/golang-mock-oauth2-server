package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/server"
)

func TestOAuthFlow(t *testing.T) {
	// Create and start the server
	mockServer := server.NewServer(":0") // Use port 0 to get a random available port

	// Start server in a goroutine
	ts := httptest.NewServer(mockServer.Handler)
	defer ts.Close()

	// Create OAuth2 config pointing to test server
	oauth2Config := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  ts.URL + "/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/authorize",
			TokenURL: ts.URL + "/token",
		},
	}

	// Test the authorization URL
	authURL := oauth2Config.AuthCodeURL("test-state", oauth2.AccessTypeOffline)

	// Create a client that doesn't follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(authURL)
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

// TestErrorScenarioIntegration tests the complete flow of:
// 1. Successful authentication
// 2. Enabling error scenario via config endpoint
// 3. Verifying authentication is blocked
// 4. Disabling error scenario
// 5. Verifying authentication succeeds again
func TestErrorScenarioIntegration(t *testing.T) {
	// Create and start the server
	mockServer := server.NewServer(":0")
	ts := httptest.NewServer(mockServer.Handler)
	defer ts.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Test parameters
	queryParams := url.Values{
		"client_id":     {"test-client"},
		"redirect_uri":  {"http://localhost/callback"},
		"scope":         {"openid"},
		"response_type": {"code"},
		"state":         {"test-state"},
	}

	t.Run("Step 1: Verify normal authentication works", func(t *testing.T) {
		authURL := ts.URL + "/authorize?" + queryParams.Encode()
		resp, err := client.Get(authURL)
		if err != nil {
			t.Fatalf("Failed to call authorize endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Errorf("Expected status %d, got %d", http.StatusFound, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		redirectURL, err := url.Parse(location)
		if err != nil {
			t.Fatalf("Failed to parse redirect URL: %v", err)
		}

		// Should have authorization code, not error
		code := redirectURL.Query().Get("code")
		if code == "" {
			t.Errorf("Expected authorization code, got none")
		}

		errorParam := redirectURL.Query().Get("error")
		if errorParam != "" {
			t.Errorf("Expected no error, got %q", errorParam)
		}
	})

	t.Run("Step 2: Enable error scenario via config endpoint", func(t *testing.T) {
		configReq := map[string]interface{}{
			"error_scenario": map[string]interface{}{
				"enabled":           true,
				"endpoint":          "authorize",
				"error":             "access_denied",
				"error_description": "User denied access via config",
			},
		}

		reqBody, err := json.Marshal(configReq)
		if err != nil {
			t.Fatalf("Failed to marshal config request: %v", err)
		}

		resp, err := http.Post(ts.URL+"/config", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to post config: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d from config endpoint, got %d", http.StatusOK, resp.StatusCode)
		}

		var configResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&configResp); err != nil {
			t.Fatalf("Failed to decode config response: %v", err)
		}

		if configResp["status"] != "success" {
			t.Errorf("Expected success status, got %v", configResp["status"])
		}
	})

	t.Run("Step 3: Verify authentication is now blocked", func(t *testing.T) {
		authURL := ts.URL + "/authorize?" + queryParams.Encode()
		resp, err := client.Get(authURL)
		if err != nil {
			t.Fatalf("Failed to call authorize endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Errorf("Expected status %d, got %d", http.StatusFound, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		redirectURL, err := url.Parse(location)
		if err != nil {
			t.Fatalf("Failed to parse redirect URL: %v", err)
		}

		// Should have error, not authorization code
		errorParam := redirectURL.Query().Get("error")
		if errorParam != "access_denied" {
			t.Errorf("Expected error 'access_denied', got %q", errorParam)
		}

		errorDesc := redirectURL.Query().Get("error_description")
		if errorDesc != "User denied access via config" {
			t.Errorf("Expected error description 'User denied access via config', got %q", errorDesc)
		}

		code := redirectURL.Query().Get("code")
		if code != "" {
			t.Errorf("Expected no authorization code when error is enabled, got %q", code)
		}
	})

	t.Run("Step 4: Disable error scenario via config endpoint", func(t *testing.T) {
		configReq := map[string]interface{}{
			"error_scenario": map[string]interface{}{
				"enabled":  false,
				"endpoint": "authorize",
			},
		}

		reqBody, err := json.Marshal(configReq)
		if err != nil {
			t.Fatalf("Failed to marshal config request: %v", err)
		}

		resp, err := http.Post(ts.URL+"/config", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to post config: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d from config endpoint, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Step 5: Verify authentication works again", func(t *testing.T) {
		authURL := ts.URL + "/authorize?" + queryParams.Encode()
		resp, err := client.Get(authURL)
		if err != nil {
			t.Fatalf("Failed to call authorize endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Errorf("Expected status %d, got %d", http.StatusFound, resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		redirectURL, err := url.Parse(location)
		if err != nil {
			t.Fatalf("Failed to parse redirect URL: %v", err)
		}

		// Should have authorization code again, not error
		code := redirectURL.Query().Get("code")
		if code == "" {
			t.Errorf("Expected authorization code after disabling error, got none")
		}

		errorParam := redirectURL.Query().Get("error")
		if errorParam != "" {
			t.Errorf("Expected no error after disabling, got %q", errorParam)
		}
	})
}

// TestTokenEndpointErrorScenario tests error scenarios for the token endpoint
func TestTokenEndpointErrorScenario(t *testing.T) {
	mockServer := server.NewServer(":0")
	ts := httptest.NewServer(mockServer.Handler)
	defer ts.Close()

	client := &http.Client{}

	// First, get an authorization code
	authClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	queryParams := url.Values{
		"client_id":     {"test-client"},
		"redirect_uri":  {"http://localhost/callback"},
		"scope":         {"openid"},
		"response_type": {"code"},
		"state":         {"test-state"},
	}

	authURL := ts.URL + "/authorize?" + queryParams.Encode()
	authResp, err := authClient.Get(authURL)
	if err != nil {
		t.Fatalf("Failed to get authorization code: %v", err)
	}
	defer authResp.Body.Close()

	location := authResp.Header.Get("Location")
	redirectURL, _ := url.Parse(location)
	authCode := redirectURL.Query().Get("code")

	if authCode == "" {
		t.Fatal("Failed to get authorization code for token test")
	}

	t.Run("Step 1: Verify token exchange works", func(t *testing.T) {
		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", authCode)
		form.Add("client_id", "test-client")
		form.Add("client_secret", "test-secret")
		form.Add("redirect_uri", "http://localhost/callback")

		resp, err := client.Post(ts.URL+"/token", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("Failed to exchange token: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var tokenResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			t.Fatalf("Failed to decode token response: %v", err)
		}

		if tokenResp["access_token"] == nil {
			t.Error("Expected access_token in response")
		}
	})

	t.Run("Step 2: Enable invalid_grant error for token endpoint", func(t *testing.T) {
		configReq := map[string]interface{}{
			"error_scenario": map[string]interface{}{
				"enabled":           true,
				"endpoint":          "token",
				"error":             "invalid_grant",
				"error_description": "Authorization code is invalid",
			},
		}

		reqBody, _ := json.Marshal(configReq)
		resp, err := http.Post(ts.URL+"/config", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to post config: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected config to succeed, got status %d", resp.StatusCode)
		}
	})

	t.Run("Step 3: Get new authorization code and verify token exchange is blocked", func(t *testing.T) {
		// Get new auth code
		authResp2, err := authClient.Get(authURL)
		if err != nil {
			t.Fatalf("Failed to get authorization code: %v", err)
		}
		defer authResp2.Body.Close()

		location2 := authResp2.Header.Get("Location")
		redirectURL2, _ := url.Parse(location2)
		newAuthCode := redirectURL2.Query().Get("code")

		// Try to exchange it
		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", newAuthCode)
		form.Add("client_id", "test-client")
		form.Add("client_secret", "test-secret")
		form.Add("redirect_uri", "http://localhost/callback")

		resp, err := client.Post(ts.URL+"/token", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("Failed to call token endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid_grant, got %d", http.StatusBadRequest, resp.StatusCode)
		}

		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if errorResp["error"] != "invalid_grant" {
			t.Errorf("Expected error 'invalid_grant', got %v", errorResp["error"])
		}

		if errorResp["error_description"] != "Authorization code is invalid" {
			t.Errorf("Expected error description 'Authorization code is invalid', got %v", errorResp["error_description"])
		}
	})
}

// TestUserInfoEndpointErrorScenario tests error scenarios for the userinfo endpoint
func TestUserInfoEndpointErrorScenario(t *testing.T) {
	mockServer := server.NewServer(":0")
	ts := httptest.NewServer(mockServer.Handler)
	defer ts.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Get an access token first
	queryParams := url.Values{
		"client_id":     {"test-client"},
		"redirect_uri":  {"http://localhost/callback"},
		"scope":         {"openid"},
		"response_type": {"code"},
		"state":         {"test-state"},
	}

	authURL := ts.URL + "/authorize?" + queryParams.Encode()
	authResp, err := client.Get(authURL)
	if err != nil {
		t.Fatalf("Failed to get authorization code: %v", err)
	}
	defer authResp.Body.Close()

	location := authResp.Header.Get("Location")
	redirectURL, _ := url.Parse(location)
	authCode := redirectURL.Query().Get("code")

	// Exchange for token
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("client_id", "test-client")
	form.Add("client_secret", "test-secret")
	form.Add("redirect_uri", "http://localhost/callback")

	tokenResp, err := http.Post(ts.URL+"/token", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	defer tokenResp.Body.Close()

	var tokenData map[string]interface{}
	json.NewDecoder(tokenResp.Body).Decode(&tokenData)
	accessToken := tokenData["access_token"].(string)

	t.Run("Step 1: Verify userinfo works with valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ts.URL+"/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to call userinfo: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Step 2: Enable server_error for userinfo endpoint", func(t *testing.T) {
		configReq := map[string]interface{}{
			"error_scenario": map[string]interface{}{
				"enabled":           true,
				"endpoint":          "userinfo",
				"error":             "server_error",
				"error_description": "Internal server error occurred",
			},
		}

		reqBody, _ := json.Marshal(configReq)
		resp, err := http.Post(ts.URL+"/config", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to post config: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected config to succeed, got status %d", resp.StatusCode)
		}
	})

	t.Run("Step 3: Verify userinfo is now blocked", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ts.URL+"/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to call userinfo: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status %d for server_error, got %d", http.StatusInternalServerError, resp.StatusCode)
		}

		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if errorResp["error"] != "server_error" {
			t.Errorf("Expected error 'server_error', got %v", errorResp["error"])
		}
	})
}
