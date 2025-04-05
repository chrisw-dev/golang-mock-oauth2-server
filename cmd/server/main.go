package main

import (
	"log"
	"net/http"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/config"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/handlers"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize in-memory store
	memoryStore := store.NewMemoryStore()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up routes
	mux.Handle("/authorize", &handlers.AuthorizeHandler{Store: memoryStore})
	mux.Handle("/token", handlers.NewTokenHandler(memoryStore))
	mux.Handle("/userinfo", &handlers.UserInfoHandler{Store: memoryStore})
	mux.Handle("/config", handlers.NewConfigHandler(memoryStore, models.NewDefaultUser()))

	// Start the server with the custom ServeMux
	startServer(cfg.Port, mux)
}

func startServer(port int, handler http.Handler) {
	addr := ":" + string(rune(port))
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
