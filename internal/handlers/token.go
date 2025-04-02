package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// TokenHandler handles OAuth2 token exchange requests
type TokenHandler struct {
	store store.Store
}

// NewTokenHandler creates a new TokenHandler with the given store
func NewTokenHandler(store store.Store) *TokenHandler {
	return &TokenHandler{
		store: store,
	}
}

// ServeHTTP handles token requests
func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Extract parameters
	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")

	// Validate grant type
	if grantType != "authorization_code" {
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
		return
	}

	// Look up authorization code
	authRequest, exists := h.store.GetAuthCode(code)
	if !exists {
		http.Error(w, "Invalid authorization code", http.StatusBadRequest)
		return
	}

	// Validate client ID
	if authRequest.ClientID != clientID {
		http.Error(w, "Client ID mismatch", http.StatusBadRequest)
		return
	}

	// Validate redirect URI
	if authRequest.RedirectURI != redirectURI {
		http.Error(w, "Redirect URI mismatch", http.StatusBadRequest)
		return
	}

	// Generate token response
	tokenResponse := models.TokenResponse{
		AccessToken:  generateAccessToken(clientID),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: generateRefreshToken(clientID),
		IDToken:      generateIDToken(clientID),
	}

	// Store the token in the store for future validation
	h.store.StoreToken(tokenResponse.AccessToken, clientID)

	// Remove the used authorization code
	h.store.RemoveAuthCode(code)

	// Return token response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse)
}

// Helper function to generate a mock access token
func generateAccessToken(clientID string) string {
	return "mock-access-token-" + clientID + "-" + time.Now().Format("20060102150405")
}

// Helper function to generate a mock refresh token
func generateRefreshToken(clientID string) string {
	return "mock-refresh-token-" + clientID + "-" + time.Now().Format("20060102150405")
}

// Helper function to generate a mock ID token
func generateIDToken(clientID string) string {
	return "mock-id-token-" + clientID + "-" + time.Now().Format("20060102150405")
}
