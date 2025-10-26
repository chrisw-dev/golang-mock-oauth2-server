package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/types"
)

// ConfigHandler handles dynamic configuration of the mock OAuth2 server
type ConfigHandler struct {
	store store.Store
	user  *models.UserInfo
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(store store.Store, defaultUser *models.UserInfo) *ConfigHandler {
	return &ConfigHandler{
		store: store,
		user:  defaultUser,
	}
}

// ConfigRequest represents a configuration request to modify server behavior
type ConfigRequest struct {
	UserInfo      map[string]interface{} `json:"user_info,omitempty"`
	Tokens        map[string]interface{} `json:"tokens,omitempty"`
	ErrorScenario *ErrorScenario         `json:"error_scenario,omitempty"`
}

// ErrorScenario defines an error condition to simulate
//
// When an error scenario is configured for an endpoint, that endpoint will return
// an OAuth2 error response instead of proceeding with the normal authentication flow.
//
// The Enabled field uses a pointer to bool (*bool) to distinguish between three states:
//   - nil (field not provided in JSON): defaults to true - error scenario is enabled
//   - true (explicitly set): error scenario is enabled
//   - false (explicitly set): error scenario is disabled
//
// This allows clients to:
//   1. Enable an error by just providing endpoint and error fields
//   2. Explicitly enable with "enabled": true
//   3. Disable a previously configured error with "enabled": false
//
// Example usage:
//   Enable error (implicit): {"endpoint": "authorize", "error": "access_denied"}
//   Enable error (explicit): {"endpoint": "authorize", "error": "access_denied", "enabled": true}
//   Disable error: {"endpoint": "authorize", "enabled": false}
type ErrorScenario struct {
	Enabled          *bool  `json:"enabled,omitempty"` // Whether the error scenario is enabled (defaults to true if not specified)
	Endpoint         string `json:"endpoint"`          // Which endpoint should return an error (authorize, token, userinfo)
	Error            string `json:"error"`             // OAuth2 error code
	ErrorDescription string `json:"error_description,omitempty"`
}

// ConfigResponse represents the response from the config endpoint
type ConfigResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ServeHTTP handles configuration requests
func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON config
	var config ConfigRequest
	err = json.Unmarshal(body, &config)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received config request: %+v", config)

	// Update user info if provided
	if config.UserInfo != nil {
		models.UpdateUserFromConfig(h.user, config.UserInfo)
	}

	// Store token configuration if provided
	if config.Tokens != nil {
		h.storeTokenConfig(config.Tokens)
	}

	// Store error scenario if provided
	if config.ErrorScenario != nil {
		h.storeErrorScenario(*config.ErrorScenario)
		log.Printf("Configured error scenario: endpoint=%s, error=%s, enabled=%v",
			config.ErrorScenario.Endpoint,
			config.ErrorScenario.Error,
			config.ErrorScenario.Enabled)
	}

	// Return success response
	response := ConfigResponse{
		Status:  "success",
		Message: "Configuration updated",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Since we've already written the header, we can only log the error
		log.Printf("Error encoding config response: %v", err)
	}
}

// storeTokenConfig saves token configuration to the store
func (h *ConfigHandler) storeTokenConfig(tokenConfig map[string]interface{}) {
	// Use the store's StoreTokenConfig method to save the token configuration
	h.store.StoreTokenConfig(tokenConfig)
}

// storeErrorScenario saves error scenario configuration to the store
//
// This function handles the defaulting logic for the Enabled field:
//   - If scenario.Enabled is nil (not provided in JSON), it defaults to true
//   - If scenario.Enabled is not nil, it uses the explicit value (true or false)
//
// The function also determines the appropriate HTTP status code based on the
// OAuth2 error code and creates a types.ErrorScenario that is stored in the store.
//
// Note: The store only maintains ONE error scenario at a time. If multiple error
// scenarios are configured for different endpoints, the last one configured will
// overwrite any previous configurations. This is by design for simplicity in a
// mock server context.
func (h *ConfigHandler) storeErrorScenario(scenario ErrorScenario) {
	// Default enabled to true when an error scenario is being configured
	// If Enabled is nil (not provided), default to true
	// If Enabled is explicitly set (true or false), use that value
	enabled := true
	if scenario.Enabled != nil {
		enabled = *scenario.Enabled
	}

	// Convert from handlers.ErrorScenario to types.ErrorScenario
	storeScenario := types.ErrorScenario{
		Enabled:     enabled,
		Endpoint:    scenario.Endpoint,
		StatusCode:  determineStatusCode(scenario.Error),
		ErrorCode:   scenario.Error,
		Description: scenario.ErrorDescription,
	}

	log.Printf("Storing error scenario: endpoint=%s, error=%s, enabled=%t, status_code=%d",
		storeScenario.Endpoint, storeScenario.ErrorCode, storeScenario.Enabled, storeScenario.StatusCode)

	// Store the error scenario in the store
	h.store.StoreErrorScenario(storeScenario)
}

// determineStatusCode returns an appropriate HTTP status code for the OAuth error
func determineStatusCode(errorCode string) int {
	switch errorCode {
	case "invalid_request":
		return http.StatusBadRequest
	case "invalid_client":
		return http.StatusUnauthorized
	case "invalid_grant":
		return http.StatusBadRequest
	case "unauthorized_client":
		return http.StatusUnauthorized
	case "unsupported_grant_type":
		return http.StatusBadRequest
	case "invalid_scope":
		return http.StatusBadRequest
	case "access_denied":
		return http.StatusForbidden
	case "unsupported_response_type":
		return http.StatusBadRequest
	case "server_error":
		return http.StatusInternalServerError
	case "temporarily_unavailable":
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadRequest
	}
}

// GetUserInfo returns the configured user info
func (h *ConfigHandler) GetUserInfo() *models.UserInfo {
	return h.user
}
