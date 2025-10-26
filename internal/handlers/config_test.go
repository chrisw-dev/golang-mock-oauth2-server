package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/types"
)

// Mock store implementation for testing
type mockStore struct {
	authCodes     map[string]*models.AuthRequest
	tokens        map[string]string
	tokenConfig   map[string]interface{}
	errorScenario *types.ErrorScenario
}

func newMockStore() *mockStore {
	return &mockStore{
		authCodes:   make(map[string]*models.AuthRequest),
		tokens:      make(map[string]string),
		tokenConfig: make(map[string]interface{}),
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}

func (s *mockStore) StoreAuthCode(code string, request *models.AuthRequest) {
	s.authCodes[code] = request
}

func (s *mockStore) GetAuthCode(code string) (*models.AuthRequest, bool) {
	req, exists := s.authCodes[code]
	return req, exists
}

func (s *mockStore) RemoveAuthCode(code string) {
	delete(s.authCodes, code)
}

func (s *mockStore) StoreToken(token string, clientID string) {
	s.tokens[token] = clientID
}

func (s *mockStore) GetClientIDByToken(token string) (string, bool) {
	clientID, exists := s.tokens[token]
	return clientID, exists
}

func (s *mockStore) StoreTokenConfig(config map[string]interface{}) {
	s.tokenConfig = config
}

func (s *mockStore) GetTokenConfig() map[string]interface{} {
	return s.tokenConfig
}

func (s *mockStore) StoreErrorScenario(scenario types.ErrorScenario) {
	s.errorScenario = &scenario
}

func (s *mockStore) GetErrorScenario(endpoint string) (*types.ErrorScenario, bool) {
	if s.errorScenario != nil && s.errorScenario.Endpoint == endpoint && s.errorScenario.Enabled {
		return s.errorScenario, true
	}
	return nil, false
}

func (s *mockStore) ClearErrorScenario(endpoint string) {
	if s.errorScenario != nil && s.errorScenario.Endpoint == endpoint {
		s.errorScenario = nil
	}
}

// Helper function to compare maps allowing numeric type differences
func compareMaps(got, want map[string]interface{}) bool {
	if len(got) != len(want) {
		return false
	}

	for k, wantVal := range want {
		gotVal, exists := got[k]
		if !exists {
			return false
		}

		// Handle numeric values specially to allow type differences
		wantNum, wantIsNum := wantVal.(float64)
		gotNum, gotIsNum := gotVal.(float64)

		// Convert int to float64 if needed
		if !gotIsNum {
			if gotInt, isInt := gotVal.(int); isInt {
				gotNum = float64(gotInt)
				gotIsNum = true
			}
		}

		if !wantIsNum {
			if wantInt, isInt := wantVal.(int); isInt {
				wantNum = float64(wantInt)
				wantIsNum = true
			}
		}

		// Compare numeric values
		if wantIsNum && gotIsNum {
			if wantNum != gotNum {
				return false
			}
		} else if !reflect.DeepEqual(gotVal, wantVal) {
			return false
		}
	}

	return true
}

func TestConfigHandler_MethodNotAllowed(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Create GET request (only POST should be allowed)
	req := httptest.NewRequest("GET", "/config", nil)
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestConfigHandler_InvalidJSON(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Create request with invalid JSON
	req := httptest.NewRequest("POST", "/config", bytes.NewBuffer([]byte(`{invalid json}`)))
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestConfigHandler_UpdateUserInfo(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Expected new values
	newName := "Updated Test User"
	newEmail := "updated@example.com"

	// Create request with user info updates
	configReq := ConfigRequest{
		UserInfo: map[string]interface{}{
			"name":  newName,
			"email": newEmail,
		},
	}
	reqBody, _ := json.Marshal(configReq)
	req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if user info was updated
	if defaultUser.Name != newName {
		t.Errorf("user name not updated: got %v want %v", defaultUser.Name, newName)
	}
	if defaultUser.Email != newEmail {
		t.Errorf("user email not updated: got %v want %v", defaultUser.Email, newEmail)
	}

	// Check response format
	var response ConfigResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if response.Status != "success" {
		t.Errorf("unexpected response status: got %v want %v", response.Status, "success")
	}
}

func TestConfigHandler_UpdateTokenConfig(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Create token config
	tokenConfig := map[string]interface{}{
		"access_token":  "custom-token",
		"expires_in":    1800,
		"refresh_token": "custom-refresh",
	}

	// Create request with token config
	configReq := ConfigRequest{
		Tokens: tokenConfig,
	}
	reqBody, _ := json.Marshal(configReq)
	req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if token config was stored
	storedConfig := mockStore.GetTokenConfig()
	if !compareMaps(storedConfig, tokenConfig) {
		t.Errorf("token config not stored correctly: got %v want %v", storedConfig, tokenConfig)
	}
}

func TestConfigHandler_UpdateErrorScenario(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Test cases for error scenarios
	testCases := []struct {
		name          string
		errorScenario ErrorScenario
		expectedCode  int
	}{
		{
			name: "invalid_request",
			errorScenario: ErrorScenario{
				Enabled:          boolPtr(true),
				Endpoint:         "token",
				Error:            "invalid_request",
				ErrorDescription: "Test invalid request",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid_client",
			errorScenario: ErrorScenario{
				Enabled:          boolPtr(true),
				Endpoint:         "token",
				Error:            "invalid_client",
				ErrorDescription: "Test invalid client",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "server_error",
			errorScenario: ErrorScenario{
				Enabled:          boolPtr(true),
				Endpoint:         "userinfo",
				Error:            "server_error",
				ErrorDescription: "Test server error",
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "unknown_error",
			errorScenario: ErrorScenario{
				Enabled:          boolPtr(true),
				Endpoint:         "authorize",
				Error:            "unknown_error",
				ErrorDescription: "Test unknown error",
			},
			expectedCode: http.StatusBadRequest, // Default for unknown errors
		},
		{
			name: "unsupported_response_type",
			errorScenario: ErrorScenario{
				Enabled:          boolPtr(true),
				Endpoint:         "authorize",
				Error:            "unsupported_response_type",
				ErrorDescription: "Response type not supported",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "temporarily_unavailable",
			errorScenario: ErrorScenario{
				Enabled:          boolPtr(true),
				Endpoint:         "authorize",
				Error:            "temporarily_unavailable",
				ErrorDescription: "Server is under maintenance",
			},
			expectedCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with error scenario
			configReq := ConfigRequest{
				ErrorScenario: &tc.errorScenario,
			}
			reqBody, _ := json.Marshal(configReq)
			req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(reqBody))
			rr := httptest.NewRecorder()

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check response code
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			// Verify error scenario was stored correctly
			scenario, exists := mockStore.GetErrorScenario(tc.errorScenario.Endpoint)
			if !exists {
				t.Errorf("error scenario not found for endpoint %s", tc.errorScenario.Endpoint)
				return
			}

			// Check the fields were stored correctly
			expectedEnabled := tc.errorScenario.Enabled != nil && *tc.errorScenario.Enabled
			if scenario.Enabled != expectedEnabled {
				t.Errorf("wrong enabled status stored: got %v want %v", scenario.Enabled, expectedEnabled)
			}
			if scenario.Endpoint != tc.errorScenario.Endpoint {
				t.Errorf("wrong endpoint stored: got %v want %v", scenario.Endpoint, tc.errorScenario.Endpoint)
			}
			if scenario.ErrorCode != tc.errorScenario.Error {
				t.Errorf("wrong error code stored: got %v want %v", scenario.ErrorCode, tc.errorScenario.Error)
			}
			if scenario.Description != tc.errorScenario.ErrorDescription {
				t.Errorf("wrong description stored: got %v want %v", scenario.Description, tc.errorScenario.ErrorDescription)
			}
			if scenario.StatusCode != tc.expectedCode {
				t.Errorf("wrong status code determined: got %v want %v", scenario.StatusCode, tc.expectedCode)
			}
		})
	}
}

func TestConfigHandler_CombinedUpdate(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Create a request that updates all three types of configuration
	configReq := ConfigRequest{
		UserInfo: map[string]interface{}{
			"name":  "Combined Test",
			"email": "combined@example.com",
		},
		Tokens: map[string]interface{}{
			"access_token": "combined-token",
			"expires_in":   2400,
		},
		ErrorScenario: &ErrorScenario{
			Enabled:          boolPtr(true),
			Endpoint:         "token",
			Error:            "invalid_grant",
			ErrorDescription: "Combined test error",
		},
	}
	reqBody, _ := json.Marshal(configReq)
	req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Assertions
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check user info was updated
	if defaultUser.Name != "Combined Test" {
		t.Errorf("user name not updated correctly")
	}
	if defaultUser.Email != "combined@example.com" {
		t.Errorf("user email not updated correctly")
	}

	// Check token config was stored
	tokenConfig := mockStore.GetTokenConfig()
	expectedTokens := map[string]interface{}{
		"access_token": "combined-token",
		"expires_in":   2400,
	}
	if !compareMaps(tokenConfig, expectedTokens) {
		t.Errorf("token config not stored correctly: got %v want %v", tokenConfig, expectedTokens)
	}

	// Check error scenario was stored
	scenario, exists := mockStore.GetErrorScenario("token")
	if !exists {
		t.Errorf("error scenario not stored")
		return
	}
	if !scenario.Enabled {
		t.Errorf("error scenario enabled flag not set correctly")
	}
	if scenario.ErrorCode != "invalid_grant" {
		t.Errorf("error scenario not stored correctly")
	}
}

func TestConfigHandler_GetUserInfo(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Get user info
	userInfo := handler.GetUserInfo()

	// Assertions
	if userInfo != defaultUser {
		t.Errorf("GetUserInfo() returned different instance than expected")
	}

	// Modify the returned user info and check if the original changed too
	userInfo.Name = "Modified Name"
	if defaultUser.Name != "Modified Name" {
		t.Errorf("UserInfo not returned by reference")
	}
}

func TestDetermineStatusCode(t *testing.T) {
	testCases := []struct {
		name      string
		errorCode string
		expected  int
	}{
		{"invalid_request", "invalid_request", http.StatusBadRequest},
		{"invalid_client", "invalid_client", http.StatusUnauthorized},
		{"invalid_grant", "invalid_grant", http.StatusBadRequest},
		{"unauthorized_client", "unauthorized_client", http.StatusUnauthorized},
		{"unsupported_grant_type", "unsupported_grant_type", http.StatusBadRequest},
		{"invalid_scope", "invalid_scope", http.StatusBadRequest},
		{"access_denied", "access_denied", http.StatusForbidden},
		{"unsupported_response_type", "unsupported_response_type", http.StatusBadRequest},
		{"server_error", "server_error", http.StatusInternalServerError},
		{"temporarily_unavailable", "temporarily_unavailable", http.StatusServiceUnavailable},
		{"unknown_error", "unknown_error", http.StatusBadRequest},
		{"empty_string", "", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := determineStatusCode(tc.errorCode)
			if result != tc.expected {
				t.Errorf("determineStatusCode(%s) = %d; want %d", tc.errorCode, result, tc.expected)
			}
		})
	}
}

func TestConfigHandler_ErrorScenarioDefaultEnabled(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Test case where enabled field is not set (nil pointer)
	// Should default to true when endpoint and error are provided
	configReq := ConfigRequest{
		ErrorScenario: &ErrorScenario{
			// Enabled field not set (nil), will default to true
			Endpoint:         "authorize",
			Error:            "unauthorized_client",
			ErrorDescription: "Client not authorized",
		},
	}
	reqBody, _ := json.Marshal(configReq)
	req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify error scenario was stored and ENABLED by default
	scenario, exists := mockStore.GetErrorScenario("authorize")
	if !exists {
		t.Errorf("error scenario not found for authorize endpoint")
		return
	}

	if !scenario.Enabled {
		t.Errorf("error scenario should be enabled by default when endpoint and error are provided, got enabled=%v", scenario.Enabled)
	}
	if scenario.ErrorCode != "unauthorized_client" {
		t.Errorf("wrong error code stored: got %v want %v", scenario.ErrorCode, "unauthorized_client")
	}
}

func TestConfigHandler_ErrorScenarioExplicitlyDisabled(t *testing.T) {
	// Setup
	mockStore := newMockStore()
	defaultUser := models.NewDefaultUser()
	handler := NewConfigHandler(mockStore, defaultUser)

	// Test case where enabled is explicitly set to false
	configReq := ConfigRequest{
		ErrorScenario: &ErrorScenario{
			Enabled:          boolPtr(false), // Explicitly disabled
			Endpoint:         "authorize",
			Error:            "access_denied",
			ErrorDescription: "Should not be active",
		},
	}
	reqBody, _ := json.Marshal(configReq)
	req := httptest.NewRequest("POST", "/config", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check response code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify error scenario is stored but NOT enabled
	scenario, exists := mockStore.GetErrorScenario("authorize")
	if exists {
		t.Errorf("error scenario should not be returned when disabled, but got: %+v", scenario)
	}
}
