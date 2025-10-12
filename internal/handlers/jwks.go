package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/jwt"
)

// JWKSHandler handles requests for JSON Web Key Set
type JWKSHandler struct{}

// NewJWKSHandler creates a new JWKS handler
func NewJWKSHandler() *JWKSHandler {
	return &JWKSHandler{}
}

// ServeHTTP handles HTTP requests for JWKS
func (h *JWKSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwks, err := jwt.GetJWKS()
	if err != nil {
		http.Error(w, "Error generating JWKS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jwks); err != nil {
		http.Error(w, "Error encoding JWKS", http.StatusInternalServerError)
		return
	}
}
