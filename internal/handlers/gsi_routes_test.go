package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterGISRoutes_EnabledVsDisabled(t *testing.T) {
	user := testGISUser()

	t.Run("disabled - routes not registered", func(t *testing.T) {
		mux := http.NewServeMux()

		for _, path := range []string{"/gsi/client", "/gsi/credential"} {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Errorf("expected 404 for %s when disabled, got %d", path, rr.Code)
			}
		}
	})

	t.Run("enabled - routes registered", func(t *testing.T) {
		mux := http.NewServeMux()
		RegisterGISRoutes(mux, user, "http://localhost:8080", nil)

		req := httptest.NewRequest(http.MethodGet, "/gsi/client", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 for /gsi/client when enabled, got %d", rr.Code)
		}

		req = httptest.NewRequest(http.MethodPost, "/gsi/credential", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Body = http.NoBody
		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400 (missing body) for /gsi/credential when enabled, got %d", rr.Code)
		}
	})
}
