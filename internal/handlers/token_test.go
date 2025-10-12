package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/jwt"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

func TestTokenHandler(t *testing.T) {
	// Setup mock store
	mockStore := store.NewMemoryStore()

	// Create test authorization code in store
	mockStore.StoreAuthCode("test-code", &models.AuthRequest{
		ClientID:    "test-client",
		RedirectURI: "http://example.com/callback",
	})

	// Create token handler with mock store
	handler := NewTokenHandler(mockStore)

	// Create test request
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", "test-code")
	form.Add("client_id", "test-client")
	form.Add("client_secret", "test-secret")
	form.Add("redirect_uri", "http://example.com/callback")

	req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check response code
	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Parse response
	var response models.TokenResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	// Check response fields
	if response.AccessToken == "" {
		t.Errorf("Expected access token to be present")
	}
	if response.TokenType != "Bearer" {
		t.Errorf("Expected token_type to be 'Bearer', got '%s'", response.TokenType)
	}

	// Verify that ID token is a valid JWT
	if response.IDToken == "" {
		t.Errorf("Expected ID token to be present")
	}

	// Parse ID token to verify it's a valid JWT
	parser := jwtlib.NewParser()
	idToken, _, err := parser.ParseUnverified(response.IDToken, jwtlib.MapClaims{})
	if err != nil {
		t.Errorf("Failed to parse ID token as JWT: %v", err)
	}

	// Check that the token has standard JWT headers
	if _, ok := idToken.Header["alg"]; !ok {
		t.Error("ID token should have 'alg' header")
	}

	if _, ok := idToken.Header["kid"]; !ok {
		t.Error("ID token should have 'kid' header")
	}

	// Verify the ID token using our JWT package
	claims, err := jwt.VerifyToken(response.IDToken)
	if err != nil {
		t.Errorf("Failed to verify ID token: %v", err)
	}

	if claims == nil {
		t.Error("ID token claims should not be nil")
	}

	// Verify that access token is a valid JWT
	accessToken, _, err := parser.ParseUnverified(response.AccessToken, jwtlib.MapClaims{})
	if err != nil {
		t.Errorf("Failed to parse access token as JWT: %v", err)
	}

	if _, ok := accessToken.Header["alg"]; !ok {
		t.Error("Access token should have 'alg' header")
	}
}
