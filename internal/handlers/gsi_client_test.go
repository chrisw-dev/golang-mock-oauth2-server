package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGISClientHandler_ServeHTTP(t *testing.T) {
	handler := NewGISClientHandler()

	req := httptest.NewRequest(http.MethodGet, "/gsi/client", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "javascript") {
		t.Errorf("expected javascript content type, got %s", contentType)
	}

	body := rr.Body.String()
	for _, marker := range []string{"g_id_onload", "data-client_id", "data-callback", "g_id_signin", "/gsi/credential"} {
		if !strings.Contains(body, marker) {
			t.Errorf("expected script body to contain %q", marker)
		}
	}
}

func TestGISClientHandler_MethodNotAllowed(t *testing.T) {
	handler := NewGISClientHandler()

	req := httptest.NewRequest(http.MethodPost, "/gsi/client", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}
