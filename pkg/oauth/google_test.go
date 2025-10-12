package oauth

import (
	"strings"
	"testing"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

func TestGoogleProvider_GenerateAuthURL(t *testing.T) {
	store := store.NewMemoryStore()
	provider := NewGoogleProvider(store)

	url := provider.GenerateAuthURL("test-client", "http://localhost/callback", "openid", "test-state")
	expected := "/authorize?client_id=test-client&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&response_type=code&scope=openid&state=test-state"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestGoogleProvider_ExchangeCodeForToken(t *testing.T) {
	store := store.NewMemoryStore()
	provider := NewGoogleProvider(store)

	// Add a valid authorization code
	code := "valid-code"
	authRequest := &models.AuthRequest{
		ClientID:   "test-client",
		Expiration: time.Now().Add(10 * time.Minute),
	}
	store.StoreAuthCode(code, authRequest)

	tests := []struct {
		name          string
		code          string
		expectedError string
	}{
		{
			name:          "Valid code",
			code:          "valid-code",
			expectedError: "",
		},
		{
			name:          "Invalid code",
			code:          "invalid-code",
			expectedError: "invalid_grant: Invalid authorization code",
		},
		{
			name:          "Expired code",
			code:          "expired-code",
			expectedError: "invalid_grant: Authorization code expired",
		},
	}

	// Add an expired authorization code
	expiredCode := "expired-code"
	expiredAuthRequest := &models.AuthRequest{
		ClientID:   "test-client",
		Expiration: time.Now().Add(-10 * time.Minute),
	}
	store.StoreAuthCode(expiredCode, expiredAuthRequest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := provider.ExchangeCodeForToken(tt.code)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("expected error %s, got %v", tt.expectedError, err)
				}
			} else {
				// Check that the result contains the expected fields
				if result == nil {
					t.Error("expected result to be non-nil")
					return
				}

				// Check token_type
				if result["token_type"] != "Bearer" {
					t.Errorf("expected token_type to be Bearer, got %v", result["token_type"])
				}

				// Check expires_in
				if result["expires_in"] != 3600 {
					t.Errorf("expected expires_in to be 3600, got %v", result["expires_in"])
				}

				// Check that access_token is a valid JWT
				accessToken, ok := result["access_token"].(string)
				if !ok || accessToken == "" {
					t.Error("access_token should be a non-empty string")
				} else {
					// Verify it's a JWT (has 3 parts separated by dots)
					parts := strings.Split(accessToken, ".")
					if len(parts) != 3 {
						t.Errorf("access_token should be a JWT with 3 parts, got %d parts", len(parts))
					}
					
					// Parse to verify it's a valid JWT
					parser := jwtlib.NewParser()
					_, _, err := parser.ParseUnverified(accessToken, jwtlib.MapClaims{})
					if err != nil {
						t.Errorf("access_token should be a valid JWT: %v", err)
					}
				}

				// Check that id_token is a valid JWT
				idToken, ok := result["id_token"].(string)
				if !ok || idToken == "" {
					t.Error("id_token should be a non-empty string")
				} else {
					// Verify it's a JWT (has 3 parts separated by dots)
					parts := strings.Split(idToken, ".")
					if len(parts) != 3 {
						t.Errorf("id_token should be a JWT with 3 parts, got %d parts", len(parts))
					}

					// Parse to verify it's a valid JWT
					parser := jwtlib.NewParser()
					_, _, err := parser.ParseUnverified(idToken, jwtlib.MapClaims{})
					if err != nil {
						t.Errorf("id_token should be a valid JWT: %v", err)
					}
				}
			}
		})
	}
}

func TestGoogleProvider_GetUserInfo(t *testing.T) {
	store := store.NewMemoryStore()
	provider := NewGoogleProvider(store)

	// Add a valid token and user info
	token := "valid-token"
	store.StoreToken(token, "test-client")
	store.StoreAuthCode("test-client", &models.AuthRequest{ClientID: "test-client"})

	tests := []struct {
		name          string
		token         string
		expectedError string
	}{
		{
			name:          "Valid token",
			token:         "valid-token",
			expectedError: "",
		},
		{
			name:          "Invalid token",
			token:         "invalid-token",
			expectedError: "invalid_token: Invalid access token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := provider.GetUserInfo(tt.token)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("expected error %s, got %v", tt.expectedError, err)
				}
			} else {
				// Check that the result contains the expected user info fields
				if result == nil {
					t.Error("expected result to be non-nil")
					return
				}

				// Check that required fields are present
				if result["sub"] == "" {
					t.Error("expected sub to be non-empty")
				}
				if result["email"] == "" {
					t.Error("expected email to be non-empty")
				}
			}
		})
	}
}
