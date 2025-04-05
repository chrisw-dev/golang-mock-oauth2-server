package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

func TestUserInfoHandler_ServeHTTP(t *testing.T) {
	store := store.NewMemoryStore()
	handler := &UserInfoHandler{Store: store}

	// Add a valid token and user info to the store
	validToken := "valid-token"
	store.StoreToken(validToken, "client-123")
	store.StoreAuthCode("client-123", &models.AuthRequest{ClientID: "client-123"})

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid token",
			authorization:  "Bearer valid-token",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"sub":"client-123","name":"Generated User","given_name":"","family_name":"","email":"client-123@example.com","email_verified":true,"picture":""}`,
		},
		{
			name:           "Missing token",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:           "Invalid token",
			authorization:  "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:           "Malformed token",
			authorization:  "InvalidHeader",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
			req.Header.Set("Authorization", tt.authorization)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}

			if tt.expectedBody != "" {
				var actualBody map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &actualBody)

				var expectedBody map[string]interface{}
				json.Unmarshal([]byte(tt.expectedBody), &expectedBody)

				if !reflect.DeepEqual(actualBody, expectedBody) {
					t.Errorf("expected body %v, got %v", expectedBody, actualBody)
				}
			}
		})
	}
}
