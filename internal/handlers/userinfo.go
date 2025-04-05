package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

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
	json.NewEncoder(w).Encode(userInfo)
}
