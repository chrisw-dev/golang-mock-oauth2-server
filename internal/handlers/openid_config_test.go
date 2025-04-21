// Package handlers provides HTTP handlers for the OAuth2 server
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenIDConfigHandler_ServeHTTP(t *testing.T) {
	// Test with different base URLs
	testCases := []struct {
		name           string
		baseURL        string
		expectedIssuer string
	}{
		{
			name:           "Local development URL",
			baseURL:        "http://localhost:8080",
			expectedIssuer: "http://localhost:8080",
		},
		{
			name:           "Production URL",
			baseURL:        "https://auth.example.com",
			expectedIssuer: "https://auth.example.com",
		},
		{
			name:           "URL with trailing slash",
			baseURL:        "https://auth.example.com/",
			expectedIssuer: "https://auth.example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewOpenIDConfigHandler(tc.baseURL)
			req := httptest.NewRequest(http.MethodGet, "/.well-known/openid-configuration", nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			// Check response status code
			if resp.Code != http.StatusOK {
				t.Errorf("expected status OK, got %v", resp.Code)
			}

			// Check content type
			contentType := resp.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			// Parse and validate the configuration
			var config map[string]interface{}
			if err := json.Unmarshal(resp.Body.Bytes(), &config); err != nil {
				t.Fatalf("failed to parse response as JSON: %v", err)
			}

			// Verify issuer
			if issuer, ok := config["issuer"].(string); !ok || issuer != tc.expectedIssuer {
				t.Errorf("expected issuer %s, got %v", tc.expectedIssuer, config["issuer"])
			}

			// Verify required endpoints are present
			requiredEndpoints := []string{
				"authorization_endpoint",
				"token_endpoint",
				"userinfo_endpoint",
				"jwks_uri",
			}

			for _, endpoint := range requiredEndpoints {
				if _, exists := config[endpoint]; !exists {
					t.Errorf("expected %s to be present in the configuration", endpoint)
				}
			}

			// Verify required arrays are present
			requiredArrays := []string{
				"response_types_supported",
				"subject_types_supported",
				"id_token_signing_alg_values_supported",
				"scopes_supported",
				"token_endpoint_auth_methods_supported",
				"claims_supported",
			}

			for _, array := range requiredArrays {
				if _, exists := config[array]; !exists {
					t.Errorf("expected %s to be present in the configuration", array)
				}
			}
		})
	}
}
