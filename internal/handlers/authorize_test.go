package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/types"
)

func TestAuthorizeHandler_ServeHTTP(t *testing.T) {
	store := store.NewMemoryStore()
	handler := &AuthorizeHandler{Store: store}

	tests := []struct {
		name           string
		queryParams    url.Values
		expectedStatus int
		expectedHeader string
	}{
		{
			name: "Valid request",
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
				"state":         {"test-state"},
			},
			expectedStatus: http.StatusFound,
			expectedHeader: "http://localhost/callback",
		},
		{
			name: "Missing client_id",
			queryParams: url.Values{
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "",
		},
		{
			name: "Invalid response_type",
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"token"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/authorize?"+tt.queryParams.Encode(), nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}

			if tt.expectedHeader != "" {
				location := resp.Header().Get("Location")
				if location == "" || location[:len(tt.expectedHeader)] != tt.expectedHeader {
					t.Errorf("expected redirect to %s, got %s", tt.expectedHeader, location)
				}
			}
		})
	}
}

func TestAuthorizeHandler_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name               string
		errorScenario      types.ErrorScenario
		queryParams        url.Values
		expectedStatus     int
		expectedError      string
		expectedErrorDesc  string
		shouldHaveState    bool
	}{
		{
			name: "access_denied error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "access_denied",
				Description: "User denied access",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
				"state":         {"test-state"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "access_denied",
			expectedErrorDesc: "User denied access",
			shouldHaveState:   true,
		},
		{
			name: "invalid_request error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "invalid_request",
				Description: "Missing client_id parameter",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "invalid_request",
			expectedErrorDesc: "Missing client_id parameter",
			shouldHaveState:   false,
		},
		{
			name: "unauthorized_client error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "unauthorized_client",
				Description: "Client not authorized",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
				"state":         {"test-state-123"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "unauthorized_client",
			expectedErrorDesc: "Client not authorized",
			shouldHaveState:   true,
		},
		{
			name: "unsupported_response_type error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "unsupported_response_type",
				Description: "Response type not supported",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "unsupported_response_type",
			expectedErrorDesc: "Response type not supported",
			shouldHaveState:   false,
		},
		{
			name: "invalid_scope error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "invalid_scope",
				Description: "Scope 'admin' is not available",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid admin"},
				"response_type": {"code"},
				"state":         {"test-state"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "invalid_scope",
			expectedErrorDesc: "Scope 'admin' is not available",
			shouldHaveState:   true,
		},
		{
			name: "server_error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "server_error",
				Description: "Internal server error",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "server_error",
			expectedErrorDesc: "Internal server error",
			shouldHaveState:   false,
		},
		{
			name: "temporarily_unavailable error",
			errorScenario: types.ErrorScenario{
				Enabled:     true,
				Endpoint:    "authorize",
				ErrorCode:   "temporarily_unavailable",
				Description: "Server is under maintenance",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
				"state":         {"test-state"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "temporarily_unavailable",
			expectedErrorDesc: "Server is under maintenance",
			shouldHaveState:   true,
		},
		{
			name: "error scenario disabled - should succeed normally",
			errorScenario: types.ErrorScenario{
				Enabled:     false, // Disabled
				Endpoint:    "authorize",
				ErrorCode:   "access_denied",
				Description: "This should not appear",
			},
			queryParams: url.Values{
				"client_id":     {"test-client"},
				"redirect_uri":  {"http://localhost/callback"},
				"scope":         {"openid"},
				"response_type": {"code"},
				"state":         {"test-state"},
			},
			expectedStatus:    http.StatusFound,
			expectedError:     "", // No error expected
			expectedErrorDesc: "",
			shouldHaveState:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh store for each test
			testStore := store.NewMemoryStore()
			handler := &AuthorizeHandler{Store: testStore}

			// Configure the error scenario
			testStore.StoreErrorScenario(tt.errorScenario)

			// Make the request
			req := httptest.NewRequest(http.MethodGet, "/authorize?"+tt.queryParams.Encode(), nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}

			// Check redirect location
			location := resp.Header().Get("Location")
			if location == "" {
				t.Fatalf("expected redirect location, got empty")
			}

			// Parse the redirect URL
			redirectURL, err := url.Parse(location)
			if err != nil {
				t.Fatalf("failed to parse redirect URL: %v", err)
			}

			// Check error parameter
			if tt.expectedError != "" {
				errorParam := redirectURL.Query().Get("error")
				if errorParam != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, errorParam)
				}

				// Check error_description parameter
				errorDesc := redirectURL.Query().Get("error_description")
				if errorDesc != tt.expectedErrorDesc {
					t.Errorf("expected error_description %q, got %q", tt.expectedErrorDesc, errorDesc)
				}
			} else {
				// When error scenario is disabled, should get authorization code
				code := redirectURL.Query().Get("code")
				if code == "" {
					t.Errorf("expected authorization code when error scenario is disabled, got none")
				}
				errorParam := redirectURL.Query().Get("error")
				if errorParam != "" {
					t.Errorf("expected no error when scenario is disabled, got %q", errorParam)
				}
			}

			// Check state parameter if expected
			if tt.shouldHaveState {
				stateParam := redirectURL.Query().Get("state")
				expectedState := tt.queryParams.Get("state")
				if stateParam != expectedState {
					t.Errorf("expected state %q, got %q", expectedState, stateParam)
				}
			}
		})
	}
}

func TestAuthorizeHandler_ErrorScenarioForDifferentEndpoint(t *testing.T) {
	// Test that error scenario for "token" endpoint doesn't affect "authorize" endpoint
	testStore := store.NewMemoryStore()
	handler := &AuthorizeHandler{Store: testStore}

	// Configure error scenario for "token" endpoint (not "authorize")
	testStore.StoreErrorScenario(types.ErrorScenario{
		Enabled:     true,
		Endpoint:    "token", // Different endpoint
		ErrorCode:   "invalid_grant",
		Description: "This should not affect authorize",
	})

	queryParams := url.Values{
		"client_id":     {"test-client"},
		"redirect_uri":  {"http://localhost/callback"},
		"scope":         {"openid"},
		"response_type": {"code"},
		"state":         {"test-state"},
	}

	req := httptest.NewRequest(http.MethodGet, "/authorize?"+queryParams.Encode(), nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	// Should succeed normally (no error)
	if resp.Code != http.StatusFound {
		t.Errorf("expected status %d, got %d", http.StatusFound, resp.Code)
	}

	location := resp.Header().Get("Location")
	redirectURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse redirect URL: %v", err)
	}

	// Should have authorization code, not error
	code := redirectURL.Query().Get("code")
	if code == "" {
		t.Errorf("expected authorization code, got none")
	}

	errorParam := redirectURL.Query().Get("error")
	if errorParam != "" {
		t.Errorf("expected no error, got %q", errorParam)
	}
}

func TestAuthorizeHandler_ErrorWithoutDescription(t *testing.T) {
	testStore := store.NewMemoryStore()
	handler := &AuthorizeHandler{Store: testStore}

	// Configure error scenario without description
	testStore.StoreErrorScenario(types.ErrorScenario{
		Enabled:     true,
		Endpoint:    "authorize",
		ErrorCode:   "access_denied",
		Description: "", // No description
	})

	queryParams := url.Values{
		"client_id":     {"test-client"},
		"redirect_uri":  {"http://localhost/callback"},
		"scope":         {"openid"},
		"response_type": {"code"},
		"state":         {"test-state"},
	}

	req := httptest.NewRequest(http.MethodGet, "/authorize?"+queryParams.Encode(), nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusFound {
		t.Errorf("expected status %d, got %d", http.StatusFound, resp.Code)
	}

	location := resp.Header().Get("Location")
	redirectURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse redirect URL: %v", err)
	}

	// Should have error but no error_description
	errorParam := redirectURL.Query().Get("error")
	if errorParam != "access_denied" {
		t.Errorf("expected error 'access_denied', got %q", errorParam)
	}

	errorDesc := redirectURL.Query().Get("error_description")
	if errorDesc != "" {
		t.Errorf("expected no error_description, got %q", errorDesc)
	}
}
