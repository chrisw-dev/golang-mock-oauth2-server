package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWKSHandler(t *testing.T) {
	handler := NewJWKSHandler()

	req := httptest.NewRequest("GET", "/jwks", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	// Parse the response
	var jwks map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&jwks)
	if err != nil {
		t.Fatalf("Failed to decode JWKS response: %v", err)
	}

	// Verify JWKS structure
	keys, ok := jwks["keys"].([]interface{})
	if !ok {
		t.Fatal("JWKS should have a 'keys' array")
	}

	if len(keys) == 0 {
		t.Error("JWKS keys array should not be empty")
	}

	// Check the first key
	key, ok := keys[0].(map[string]interface{})
	if !ok {
		t.Fatal("Key should be a map")
	}

	requiredFields := []string{"kty", "use", "kid", "alg", "n", "e"}
	for _, field := range requiredFields {
		if _, exists := key[field]; !exists {
			t.Errorf("Key should have field %s", field)
		}
	}

	// Verify key type is RSA
	if key["kty"] != "RSA" {
		t.Errorf("Expected kty to be RSA, got %v", key["kty"])
	}

	// Verify algorithm is RS256
	if key["alg"] != "RS256" {
		t.Errorf("Expected alg to be RS256, got %v", key["alg"])
	}

	// Verify use is sig
	if key["use"] != "sig" {
		t.Errorf("Expected use to be sig, got %v", key["use"])
	}
}
