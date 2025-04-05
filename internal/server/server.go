package server

import (
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/handlers"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

type Server struct {
	Handler http.Handler
}

func NewServer(addr string) *Server {
	store := store.NewMemoryStore()
	authorizeHandler := &handlers.AuthorizeHandler{Store: store}
	tokenHandler := handlers.NewTokenHandler(store)
	callbackHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux := http.NewServeMux()
	mux.Handle("/authorize", authorizeHandler)
	mux.Handle("/token", tokenHandler)
	mux.Handle("/callback", callbackHandler)

	return &Server{
		Handler: mux,
	}
}
