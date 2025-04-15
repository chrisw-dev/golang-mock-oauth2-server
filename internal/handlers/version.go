// Package handlers provides HTTP handlers for the OAuth2 server
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/version"
)

// VersionHandler provides version information about the server
type VersionHandler struct{}

// NewVersionHandler creates a new VersionHandler
func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

// ServeHTTP handles HTTP requests for version information
func (h *VersionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info := version.GetVersion()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Error encoding version response", http.StatusInternalServerError)
		return
	}
}
