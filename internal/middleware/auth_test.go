package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

func TestValidateToken(t *testing.T) {
	store := store.NewMemoryStore()
	middleware := NewAuthMiddleware(store)

	// Add a valid token to the store
	validToken := "valid-token"
	store.StoreToken(validToken, "client-123")
	store.StoreAuthCode("client-123", &models.AuthRequest{ClientID: "client-123"})

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			authorization:  "Bearer valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing token",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			authorization:  "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed token",
			authorization:  "InvalidHeader",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.authorization)
			resp := httptest.NewRecorder()

			handler := middleware.ValidateToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
		})
	}
}
