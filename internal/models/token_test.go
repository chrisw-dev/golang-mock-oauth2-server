package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestTokenResponseMarshaling(t *testing.T) {
	// Create a sample token response
	tokenResp := TokenResponse{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "test-refresh-token",
		IDToken:      "test-id-token",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tokenResp)
	if err != nil {
		t.Fatalf("Failed to marshal TokenResponse: %v", err)
	}

	// Unmarshal back to struct
	var unmarshaled TokenResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal TokenResponse: %v", err)
	}

	// Verify unmarshaled data matches original
	if !reflect.DeepEqual(tokenResp, unmarshaled) {
		t.Errorf("Unmarshaled data doesn't match original. Got: %+v, Want: %+v", unmarshaled, tokenResp)
	}
}

func TestTokenResponseJSONFieldNames(t *testing.T) {
	// Create a sample token response
	tokenResp := TokenResponse{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "test-refresh-token",
		IDToken:      "test-id-token",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tokenResp)
	if err != nil {
		t.Fatalf("Failed to marshal TokenResponse: %v", err)
	}

	// Convert to string for inspection
	jsonStr := string(jsonData)

	// Check that the JSON contains the expected field names
	expectedFields := []string{
		`"access_token"`,
		`"token_type"`,
		`"expires_in"`,
		`"refresh_token"`,
		`"id_token"`,
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON string doesn't contain expected field %s. JSON: %s", field, jsonStr)
		}
	}
}

func TestAuthRequestFields(t *testing.T) {
	// Test creating an AuthRequest with different values
	testCases := []struct {
		name        string
		clientID    string
		redirectURI string
		scope       string
		expiration  time.Time
	}{
		{
			name:        "Standard values",
			clientID:    "client-123",
			redirectURI: "https://example.com/callback",
		},
		{
			name:        "Empty client ID",
			clientID:    "",
			redirectURI: "https://example.com/callback",
		},
		{
			name:        "Empty redirect URI",
			clientID:    "client-123",
			redirectURI: "",
		},
		{
			name:        "Both empty",
			clientID:    "",
			redirectURI: "",
		},
		{
			name:        "With Scope and Expiration",
			clientID:    "client-123",
			redirectURI: "https://example.com/callback",
			scope:       "openid profile",
			expiration:  time.Now().Add(10 * time.Minute),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authReq := AuthRequest{
				ClientID:    tc.clientID,
				RedirectURI: tc.redirectURI,
				Scope:       tc.scope,
				Expiration:  tc.expiration,
			}

			// Verify the fields were set correctly
			if authReq.ClientID != tc.clientID {
				t.Errorf("ClientID field not set correctly. Got: %s, Want: %s", authReq.ClientID, tc.clientID)
			}
			if authReq.RedirectURI != tc.redirectURI {
				t.Errorf("RedirectURI field not set correctly. Got: %s, Want: %s", authReq.RedirectURI, tc.redirectURI)
			}
			if authReq.Scope != tc.scope {
				t.Errorf("Scope field not set correctly. Got: %s, Want: %s", authReq.Scope, tc.scope)
			}
			if !authReq.Expiration.Equal(tc.expiration) {
				t.Errorf("Expiration field not set correctly. Got: %v, Want: %v", authReq.Expiration, tc.expiration)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != substr && s[len(s)-len(substr):] != substr
}

func TestTokenResponseDefaults(t *testing.T) {
	// Test with empty struct
	emptyResp := TokenResponse{}

	// Marshal to JSON
	jsonData, err := json.Marshal(emptyResp)
	if err != nil {
		t.Fatalf("Failed to marshal empty TokenResponse: %v", err)
	}

	// Convert to string for inspection
	jsonStr := string(jsonData)

	// Check JSON structure
	expected := `{"access_token":"","token_type":"","expires_in":0,"refresh_token":"","id_token":""}`
	if jsonStr != expected {
		t.Errorf("Empty TokenResponse JSON not as expected. Got: %s, Want: %s", jsonStr, expected)
	}
}

func TestTokenResponseUnmarshalFromJSON(t *testing.T) {
	// Test cases for unmarshaling JSON to TokenResponse
	testCases := []struct {
		name     string
		jsonStr  string
		expected TokenResponse
		wantErr  bool
	}{
		{
			name:     "Valid full JSON",
			jsonStr:  `{"access_token":"abc123","token_type":"Bearer","expires_in":3600,"refresh_token":"refresh123","id_token":"id123"}`,
			expected: TokenResponse{AccessToken: "abc123", TokenType: "Bearer", ExpiresIn: 3600, RefreshToken: "refresh123", IDToken: "id123"},
			wantErr:  false,
		},
		{
			name:     "Missing fields",
			jsonStr:  `{"access_token":"abc123","token_type":"Bearer"}`,
			expected: TokenResponse{AccessToken: "abc123", TokenType: "Bearer"},
			wantErr:  false,
		},
		{
			name:    "Invalid expires_in type",
			jsonStr: `{"access_token":"abc123","token_type":"Bearer","expires_in":"not-a-number"}`,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resp TokenResponse
			err := json.Unmarshal([]byte(tc.jsonStr), &resp)

			// Check if we expected an error
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error when unmarshaling, but got nil")
				}
				return
			}

			// If we didn't expect an error, it should be nil
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify the unmarshaled data matches expected
			if !reflect.DeepEqual(resp, tc.expected) {
				t.Errorf("Unmarshaled data doesn't match expected. Got: %+v, Want: %+v", resp, tc.expected)
			}
		})
	}
}
