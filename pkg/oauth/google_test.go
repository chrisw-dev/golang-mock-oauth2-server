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

func TestGoogleProvider_ExchangeCodeForToken_WithConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		tokenConfig   map[string]interface{}
		expectedEmail string
		expectedName  string
	}{
		{
			name:          "No configuration - uses defaults",
			tokenConfig:   nil,
			expectedEmail: "test-client@example.com",
			expectedName:  "User test-client",
		},
		{
			name: "Full configuration - uses configured values",
			tokenConfig: map[string]interface{}{
				"user_info": map[string]interface{}{
					"email": "custom@example.com",
					"name":  "Custom User",
				},
			},
			expectedEmail: "custom@example.com",
			expectedName:  "Custom User",
		},
		{
			name: "Partial configuration - email only",
			tokenConfig: map[string]interface{}{
				"user_info": map[string]interface{}{
					"email": "partial@example.com",
				},
			},
			expectedEmail: "partial@example.com",
			expectedName:  "User test-client",
		},
		{
			name: "Partial configuration - name only",
			tokenConfig: map[string]interface{}{
				"user_info": map[string]interface{}{
					"name": "Partial User",
				},
			},
			expectedEmail: "test-client@example.com",
			expectedName:  "Partial User",
		},
		{
			name: "Empty configuration object",
			tokenConfig: map[string]interface{}{
				"user_info": map[string]interface{}{},
			},
			expectedEmail: "test-client@example.com",
			expectedName:  "User test-client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := store.NewMemoryStore()
			provider := NewGoogleProvider(store)

			// Set up token configuration if provided
			if tt.tokenConfig != nil {
				store.StoreTokenConfig(tt.tokenConfig)
			}

			// Add a valid authorization code
			code := "valid-code"
			authRequest := &models.AuthRequest{
				ClientID:   "test-client",
				Expiration: time.Now().Add(10 * time.Minute),
			}
			store.StoreAuthCode(code, authRequest)

			// Exchange the code for a token
			result, err := provider.ExchangeCodeForToken(code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Get the ID token
			idToken, ok := result["id_token"].(string)
			if !ok || idToken == "" {
				t.Fatal("id_token should be a non-empty string")
			}

			// Parse the ID token to verify claims
			parser := jwtlib.NewParser()
			token, _, err := parser.ParseUnverified(idToken, jwtlib.MapClaims{})
			if err != nil {
				t.Fatalf("failed to parse ID token: %v", err)
			}

			claims, ok := token.Claims.(jwtlib.MapClaims)
			if !ok {
				t.Fatal("failed to get claims from token")
			}

			// Verify email claim
			email, hasEmail := claims["email"].(string)
			if tt.expectedEmail != "" {
				if !hasEmail {
					t.Errorf("expected email claim to be present")
				} else if email != tt.expectedEmail {
					t.Errorf("expected email=%s, got %s", tt.expectedEmail, email)
				}
			}

			// Verify name claim
			name, hasName := claims["name"].(string)
			if tt.expectedName != "" {
				if !hasName {
					t.Errorf("expected name claim to be present")
				} else if name != tt.expectedName {
					t.Errorf("expected name=%s, got %s", tt.expectedName, name)
				}
			}

			// Verify other expected claims are present
			if claims["sub"] == "" {
				t.Error("expected sub claim to be present")
			}
			if claims["aud"] != "test-client" {
				t.Errorf("expected aud=test-client, got %v", claims["aud"])
			}
			if claims["iss"] == "" {
				t.Error("expected iss claim to be present")
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
