package handlers

import (
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
)

// RegisterGISRoutes registers the opt-in mock Google Identity Services routes
// (/gsi/client and /gsi/credential) on the given mux. Callers must only invoke
// this when GIS compatibility is explicitly enabled.
func RegisterGISRoutes(mux *http.ServeMux, user *models.UserInfo, issuerURL string, allowedOrigins []string) {
	mux.Handle("/gsi/client", NewGISClientHandler())
	mux.Handle("/gsi/credential", NewGISCredentialHandler(user, issuerURL, allowedOrigins))
}
