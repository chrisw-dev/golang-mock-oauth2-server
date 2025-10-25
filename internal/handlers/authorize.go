package handlers

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"

	"github.com/google/uuid"
)

// AuthorizeHandler handles OAuth2 authorization requests
type AuthorizeHandler struct {
	Store *store.MemoryStore
}

func (h *AuthorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	responseType := r.URL.Query().Get("response_type")
	state := r.URL.Query().Get("state")

	if clientID == "" || redirectURI == "" || scope == "" || responseType != "code" {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	// Check if there's an error scenario configured for the authorize endpoint
	if errorScenario, exists := h.Store.GetErrorScenario("authorize"); exists {
		// Redirect to the provided redirect URI with error parameters
		redirectURL, err := url.Parse(redirectURI)
		if err != nil {
			http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
			return
		}

		query := redirectURL.Query()
		query.Set("error", errorScenario.ErrorCode)
		if errorScenario.Description != "" {
			query.Set("error_description", errorScenario.Description)
		}
		if state != "" {
			query.Set("state", state)
		}
		redirectURL.RawQuery = query.Encode()

		log.Printf("Returning error redirect for authorize endpoint: error=%s, description=%s", errorScenario.ErrorCode, errorScenario.Description)
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		return
	}

	// Generate authorization code
	authCode := uuid.New().String()
	expiration := time.Now().Add(10 * time.Minute)

	// Store the authorization code
	h.Store.StoreAuthCode(authCode, &models.AuthRequest{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		Expiration:  expiration,
	})

	// Redirect to the provided redirect URI with the authorization code
	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
		return
	}

	query := redirectURL.Query()
	query.Set("code", authCode)
	if state != "" {
		query.Set("state", state)
	}
	redirectURL.RawQuery = query.Encode()

	// Added logging to debug query parameters and response status.
	log.Printf("Received request with query parameters: client_id=%s, redirect_uri=%s, scope=%s, response_type=%s, state=%s", clientID, redirectURI, scope, responseType, state)
	log.Printf("Returning redirect to: %s", redirectURL.String())

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
