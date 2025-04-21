// Package handlers provides HTTP handlers for the OAuth2 server
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// OpenIDConfigHandler handles requests to the OpenID Connect discovery endpoint
type OpenIDConfigHandler struct {
	BaseURL string
}

// NewOpenIDConfigHandler creates a new OpenID Connect configuration handler
func NewOpenIDConfigHandler(baseURL string) *OpenIDConfigHandler {
	// Ensure the baseURL doesn't end with a slash
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &OpenIDConfigHandler{
		BaseURL: baseURL,
	}
}

// ServeHTTP handles HTTP requests for OpenID Connect configuration
func (h *OpenIDConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create the OpenID Connect configuration
	config := map[string]interface{}{
		"issuer":                                h.BaseURL,
		"authorization_endpoint":                h.BaseURL + "/authorize",
		"token_endpoint":                        h.BaseURL + "/token",
		"userinfo_endpoint":                     h.BaseURL + "/userinfo",
		"jwks_uri":                              h.BaseURL + "/jwks",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "email", "profile"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_post", "client_secret_basic"},
		"claims_supported": []string{
			"sub",
			"iss",
			"name",
			"given_name",
			"family_name",
			"email",
			"email_verified",
			"picture",
		},
	}

	// Return the configuration as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "Error encoding OpenID configuration", http.StatusInternalServerError)
		return
	}
}
