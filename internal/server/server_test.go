package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer(":8080")

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Valid /authorize endpoint",
			path:           "/authorize",
			expectedStatus: http.StatusBadRequest, // Expecting BadRequest due to missing query params
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

			server.Handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
