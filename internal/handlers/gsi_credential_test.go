package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/jwt"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

func testGISUser() *models.UserInfo {
	return &models.UserInfo{
		Sub:           "123456789",
		Name:          "Test User",
		Email:         "testuser@example.com",
		EmailVerified: true,
		Picture:       "https://example.com/profile.jpg",
	}
}

func TestGISCredentialHandler_Success(t *testing.T) {
	handler := NewGISCredentialHandler(testGISUser(), "http://localhost:8080", nil)

	body, _ := json.Marshal(map[string]string{"client_id": "demo-local"})
	req := httptest.NewRequest(http.MethodPost, "/gsi/credential", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp gisCredentialResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Credential == "" {
		t.Fatal("expected non-empty credential")
	}

	claims, err := jwt.VerifyToken(resp.Credential)
	if err != nil {
		t.Fatalf("failed to verify credential: %v", err)
	}
	if claims["iss"] != "http://localhost:8080" {
		t.Errorf("expected issuer http://localhost:8080, got %v", claims["iss"])
	}
	if claims["aud"] != "demo-local" {
		t.Errorf("expected audience demo-local, got %v", claims["aud"])
	}
	if claims["email"] != "testuser@example.com" {
		t.Errorf("expected email testuser@example.com, got %v", claims["email"])
	}
	if emailVerified, ok := claims["email_verified"].(bool); !ok || !emailVerified {
		t.Errorf("expected email_verified true, got %v", claims["email_verified"])
	}
}

func TestGISCredentialHandler_MethodNotAllowed(t *testing.T) {
	handler := NewGISCredentialHandler(testGISUser(), "http://localhost:8080", nil)

	req := httptest.NewRequest(http.MethodGet, "/gsi/credential", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestGISCredentialHandler_MissingClientID(t *testing.T) {
	handler := NewGISCredentialHandler(testGISUser(), "http://localhost:8080", nil)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/gsi/credential", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGISCredentialHandler_MalformedJSON(t *testing.T) {
	handler := NewGISCredentialHandler(testGISUser(), "http://localhost:8080", nil)

	req := httptest.NewRequest(http.MethodPost, "/gsi/credential", bytes.NewReader([]byte("{not-json")))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGISCredentialHandler_OriginChecks(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		origin         string
		allowedOrigins []string
		expectedStatus int
	}{
		{
			name:           "no origin header allowed",
			host:           "localhost:8080",
			origin:         "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "matching same-origin allowed",
			host:           "localhost:8080",
			origin:         "http://localhost:8080",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "cross-origin rejected by default",
			host:           "localhost:8080",
			origin:         "http://evil.example.com",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "matching host but mismatched scheme rejected",
			host:           "localhost:8080",
			origin:         "https://localhost:8080",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "explicitly allowed origin accepted",
			host:           "localhost:8080",
			origin:         "http://localhost:5173",
			allowedOrigins: []string{"http://localhost:5173"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGISCredentialHandler(testGISUser(), "http://localhost:8080", tt.allowedOrigins)

			body, _ := json.Marshal(map[string]string{"client_id": "demo-local"})
			req := httptest.NewRequest(http.MethodPost, "/gsi/credential", bytes.NewReader(body))
			req.Host = tt.host
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestGISCredentialHandler_ConcurrentWithConfigUpdates guards against data races between
// /config mutating the shared mock user in place and /gsi/credential reading it concurrently.
func TestGISCredentialHandler_ConcurrentWithConfigUpdates(t *testing.T) {
	user := testGISUser()
	mockStore := store.NewMemoryStore()
	configHandler := NewConfigHandler(mockStore, user)
	credentialHandler := NewGISCredentialHandler(user, "http://localhost:8080", nil)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			body, _ := json.Marshal(map[string]interface{}{
				"user_info": map[string]interface{}{"name": "Concurrent User"},
			})
			req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			configHandler.ServeHTTP(rr, req)
		}()

		go func() {
			defer wg.Done()
			body, _ := json.Marshal(map[string]string{"client_id": "demo-local"})
			req := httptest.NewRequest(http.MethodPost, "/gsi/credential", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			credentialHandler.ServeHTTP(rr, req)
		}()
	}
	wg.Wait()
}

