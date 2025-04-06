package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// UserInfoHandler handles requests to the OAuth2 userinfo endpoint
type UserInfoHandler struct {
	Store *store.MemoryStore
}

func (h *UserInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := authHeader[7:]
	userInfo, exists := h.Store.GetUserInfoByToken(token)
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// Log the error for debugging purposes
		log.Printf("Error encoding user info response: %v", err)
		return
	}
}
