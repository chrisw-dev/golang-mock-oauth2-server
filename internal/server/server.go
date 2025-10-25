package server

import (
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/handlers"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// Server represents an OAuth2 mock server
type Server struct {
	Handler http.Handler
}

// NewServer creates a new OAuth2 server instance with configured routes
func NewServer(addr string) *Server {
	// Initialize in-memory store
	memoryStore := store.NewMemoryStore()
	
	// Set up default user
	defaultUser := models.NewDefaultUser()
	
	// Create handlers
	authorizeHandler := &handlers.AuthorizeHandler{Store: memoryStore}
	tokenHandler := handlers.NewTokenHandler(memoryStore)
	userInfoHandler := &handlers.UserInfoHandler{Store: memoryStore}
	configHandler := handlers.NewConfigHandler(memoryStore, defaultUser)
	versionHandler := handlers.NewVersionHandler()
	jwksHandler := handlers.NewJWKSHandler()
	openIDConfigHandler := handlers.NewOpenIDConfigHandler("http://localhost" + addr)
	
	callbackHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux := http.NewServeMux()
	mux.Handle("/authorize", authorizeHandler)
	mux.Handle("/token", tokenHandler)
	mux.Handle("/userinfo", userInfoHandler)
	mux.Handle("/config", configHandler)
	mux.Handle("/version", versionHandler)
	mux.Handle("/jwks", jwksHandler)
	mux.Handle("/.well-known/openid-configuration", openIDConfigHandler)
	mux.Handle("/callback", callbackHandler)

	return &Server{
		Handler: mux,
	}
}
