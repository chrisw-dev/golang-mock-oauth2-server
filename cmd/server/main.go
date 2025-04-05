package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

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
	addr := ":" + strconv.Itoa(port)
	log.Printf("Starting server on %s...", addr)

	// Create a server with proper timeout settings
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start the server with graceful shutdown capabilities
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
