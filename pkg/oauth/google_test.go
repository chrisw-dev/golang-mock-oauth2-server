package oauth

import (
	"reflect"
	"testing"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
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
		name           string
		code           string
		expectedError  string
		expectedResult map[string]interface{}
	}{
		{
			name:          "Valid code",
			code:          "valid-code",
			expectedError: "",
			expectedResult: map[string]interface{}{
				"access_token":  "mock-access-token",
				"token_type":    "Bearer",
				"expires_in":    3600,
				"refresh_token": "mock-refresh-token",
				"id_token":      "mock-id-token",
			},
		},
		{
			name:           "Invalid code",
			code:           "invalid-code",
			expectedError:  "invalid_grant: Invalid authorization code",
			expectedResult: nil,
		},
		{
			name:           "Expired code",
			code:           "expired-code",
			expectedError:  "invalid_grant: Authorization code expired",
			expectedResult: nil,
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
			} else if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("expected result %+v, got %+v", tt.expectedResult, result)
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
		name           string
		token          string
		expectedError  string
		expectedResult map[string]interface{}
	}{
		{
			name:          "Valid token",
			token:         "valid-token",
			expectedError: "",
			expectedResult: map[string]interface{}{
				"sub":            "test-client",
				"name":           "Generated User",
				"given_name":     "",
				"family_name":    "",
				"email":          "test-client@example.com",
				"email_verified": true,
				"picture":        "",
				"locale":         "",
				"hd":             "",
			},
		},
		{
			name:           "Invalid token",
			token:          "invalid-token",
			expectedError:  "invalid_token: Invalid access token",
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := provider.GetUserInfo(tt.token)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("expected error %s, got %v", tt.expectedError, err)
				}
			} else if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("expected result %+v, got %+v", tt.expectedResult, result)
			}
		})
	}
}
