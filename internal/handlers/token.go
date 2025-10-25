package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/jwt"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// TokenHandler handles OAuth2 token exchange requests
type TokenHandler struct {
	store     store.Store
	issuerURL string
}

// NewTokenHandler creates a new TokenHandler with the given store
func NewTokenHandler(store store.Store) *TokenHandler {
	return &TokenHandler{
		store:     store,
		issuerURL: "http://localhost:8080", // default issuer
	}
}

// NewTokenHandlerWithIssuer creates a new TokenHandler with the given store and issuer URL
func NewTokenHandlerWithIssuer(store store.Store, issuerURL string) *TokenHandler {
	return &TokenHandler{
		store:     store,
		issuerURL: issuerURL,
	}
}

// ServeHTTP handles token requests
func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for error scenarios configured for the token endpoint
	if errorScenario, exists := h.store.GetErrorScenario("token"); exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorScenario.StatusCode)
		
		errorResponse := map[string]string{
			"error": errorScenario.ErrorCode,
		}
		
		if errorScenario.Description != "" {
			errorResponse["error_description"] = errorScenario.Description
		}
		
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Printf("Error encoding error response: %v", err)
		}
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
	accessToken, err := generateAccessToken(h.issuerURL, clientID, authRequest.Scope)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	idToken, err := h.generateIDToken(h.issuerURL, clientID)
	if err != nil {
		log.Printf("Error generating ID token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenResponse := models.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: generateRefreshToken(clientID),
		IDToken:      idToken,
	}

	// Store the token in the store for future validation
	h.store.StoreToken(tokenResponse.AccessToken, clientID)

	// Remove the used authorization code
	h.store.RemoveAuthCode(code)

	// Return token response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// Log the error for debugging purposes
		log.Printf("Error encoding token response: %v", err)
		return
	}
}

// Helper function to generate a mock access token
func generateAccessToken(issuerURL, clientID, scope string) (string, error) {
	// Parse scopes from the scope string
	scopes := strings.Fields(scope)
	if len(scopes) == 0 {
		scopes = []string{"openid"}
	}

	// Generate a subject ID based on client ID
	sub := "user-" + clientID

	return jwt.GenerateAccessToken(issuerURL, clientID, sub, scopes)
}

// Helper function to generate a mock refresh token
func generateRefreshToken(clientID string) string {
	return "mock-refresh-token-" + clientID + "-" + time.Now().Format("20060102150405")
}

// Helper function to generate a mock ID token
func (h *TokenHandler) generateIDToken(issuerURL, clientID string) (string, error) {
	// Generate a subject ID based on client ID
	sub := "user-" + clientID

	// Check if there's a configured email in the token config
	var email string
	var name string
	tokenConfig := h.store.GetTokenConfig()
	if tokenConfig != nil {
		if userInfoConfig, ok := tokenConfig["user_info"].(map[string]interface{}); ok {
			if configuredEmail, ok := userInfoConfig["email"].(string); ok {
				email = configuredEmail
			}
			if configuredName, ok := userInfoConfig["name"].(string); ok {
				name = configuredName
			}
		}
	}

	// If no email is configured, pass empty string (don't default to generated email)
	return jwt.GenerateIDToken(issuerURL, clientID, sub, email, name)
}
