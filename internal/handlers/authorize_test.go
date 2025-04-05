package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
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
