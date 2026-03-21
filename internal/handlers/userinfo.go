package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// UserInfoHandler handles requests to the OAuth2 userinfo endpoint
type UserInfoHandler struct {
	Store *store.MemoryStore
}

func (h *UserInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log request details
	log.Printf("UserInfo request received from %s", sanitizeLog(r.RemoteAddr)) // #nosec G706 -- sanitizeLog strips newlines/CRs to prevent log injection

	// Check for error scenarios configured for the userinfo endpoint
	if errorScenario, exists := h.Store.GetErrorScenario("userinfo"); exists {
		log.Printf("UserInfo request: Returning configured error: %s", errorScenario.ErrorCode)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorScenario.StatusCode)
		
		errorResponse := map[string]string{
		"error": errorScenario.ErrorCode,
	}
	
	if errorScenario.Description != "" {
		errorResponse["error_description"] = errorScenario.Description
	}
	
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
	return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("UserInfo request failed: Authorization header missing")
		http.Error(w, "Unauthorized - Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Printf("UserInfo request failed: Invalid Authorization header format: %s", sanitizeLog(authHeader)) // #nosec G706 -- sanitizeLog strips newlines/CRs to prevent log injection
		http.Error(w, "Unauthorized - Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	token := authHeader[7:]
	log.Printf("UserInfo request: Validating token: %s", sanitizeLog(maskToken(token))) // #nosec G706 -- sanitizeLog strips newlines/CRs to prevent log injection

	userInfo, exists := h.Store.GetUserInfoByToken(token)
	if !exists {
		log.Printf("UserInfo request failed: Token not found or invalid")
		http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("UserInfo request successful for user: %s", sanitizeLog(userInfo.Email)) // #nosec G706 -- sanitizeLog strips newlines/CRs to prevent log injection
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		log.Printf("Error encoding user info response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// maskToken hides most of the token for security in logs
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
