package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	ContextKeyUserInfo ContextKey = "userInfo"
)

type AuthMiddleware struct {
	Store *store.MemoryStore
}

func NewAuthMiddleware(store *store.MemoryStore) *AuthMiddleware {
	return &AuthMiddleware{Store: store}
}

func (a *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userInfo, exists := a.Store.GetUserInfoByToken(token)
		if !exists {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to the request context
		ctx := context.WithValue(r.Context(), ContextKeyUserInfo, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
