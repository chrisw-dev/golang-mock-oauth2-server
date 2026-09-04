package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/jwt"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
)

// GISCredentialHandler issues a mock Google Identity Services credential (RS256 ID token)
// for the active mock user, for local development/test use only.
type GISCredentialHandler struct {
	User           *models.UserInfo
	IssuerURL      string
	AllowedOrigins []string
}

// NewGISCredentialHandler creates a new GISCredentialHandler
func NewGISCredentialHandler(user *models.UserInfo, issuerURL string, allowedOrigins []string) *GISCredentialHandler {
	return &GISCredentialHandler{
		User:           user,
		IssuerURL:      issuerURL,
		AllowedOrigins: allowedOrigins,
	}
}

// gisCredentialRequest represents the JSON body of a credential request
type gisCredentialRequest struct {
	ClientID string `json:"client_id"`
}

// gisCredentialResponse represents the JSON response returned to the GIS client script
type gisCredentialResponse struct {
	Credential string `json:"credential"`
}

// ServeHTTP handles requests for mock GIS credentials
func (h *GISCredentialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	origin := r.Header.Get("Origin")
	if !h.isAllowedOrigin(r) {
		log.Printf("GSI credential request rejected: disallowed origin %s", sanitizeLog(origin)) // #nosec G706 -- sanitizeLog strips CR/LF to prevent log injection
		http.Error(w, "Forbidden - Origin not allowed", http.StatusForbidden)
		return
	}

	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Add("Vary", "Origin")
	}
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", http.MethodPost)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limit request body to 1MB
	var req gisCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ClientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}

	// Snapshot the shared mock user under lock; /config may mutate it concurrently.
	userMu.RLock()
	user := h.User.Clone()
	userMu.RUnlock()

	credential, err := jwt.GenerateGISCredentialToken(h.IssuerURL, req.ClientID, user)
	if err != nil {
		log.Printf("Error generating GIS credential: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gisCredentialResponse{Credential: credential}); err != nil {
		log.Printf("Error encoding GIS credential response: %v", err)
	}
}

// isAllowedOrigin enforces same-origin access by default: requests without an Origin header
// (e.g. same-origin browser requests, curl, tests) are allowed; cross-origin requests are
// only allowed if they match the request's full scheme+host origin or an explicitly
// configured allowed origin.
func (h *GISCredentialHandler) isAllowedOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	if origin == requestOrigin(r) {
		return true
	}

	for _, allowed := range h.AllowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

// requestOrigin reconstructs the scheme://host origin of the incoming request for
// comparison against the browser-supplied Origin header (which includes the scheme).
func requestOrigin(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
