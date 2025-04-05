package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/handlers"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// Updated to use the same ServeMux instance as in `main.go`.
func TestMainRoutes(t *testing.T) {
	mux := http.NewServeMux()

	// Initialize in-memory store
	memoryStore := store.NewMemoryStore()

	// Set up routes
	mux.Handle("/authorize", &handlers.AuthorizeHandler{Store: memoryStore})
	mux.Handle("/token", handlers.NewTokenHandler(memoryStore))
	mux.Handle("/userinfo", &handlers.UserInfoHandler{Store: memoryStore})
	mux.Handle("/config", handlers.NewConfigHandler(memoryStore, models.NewDefaultUser()))

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Authorize endpoint",
			path:           "/authorize",
			expectedStatus: http.StatusBadRequest, // Missing query params
		},
		{
			name:           "Token endpoint",
			path:           "/token",
			expectedStatus: http.StatusMethodNotAllowed, // Expecting POST
		},
		{
			name:           "Userinfo endpoint",
			path:           "/userinfo",
			expectedStatus: http.StatusUnauthorized, // Missing Authorization header
		},
		{
			name:           "Config endpoint",
			path:           "/config",
			expectedStatus: http.StatusMethodNotAllowed, // Expecting POST
		},
		{
			name:           "Invalid endpoint",
			path:           "/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp := httptest.NewRecorder()

			mux.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
