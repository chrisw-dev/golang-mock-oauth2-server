package main

import (
	"flag"
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
	// Define command-line flags
	var port int
	flag.IntVar(&port, "port", 0, "Port to run the server on (default: uses MOCK_OAUTH_PORT env var or 8080)")
	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	// Default port is 8080
	serverPort := 8080

	// Check environment variable first (using config to access env vars consistently)
	if cfg.Port > 0 {
		serverPort = cfg.Port
	}

	// Command-line flag overrides environment variable
	if port > 0 {
		serverPort = port
	}

	// Initialize in-memory store with configuration
	memoryStore := store.NewMemoryStore()

	// Set up default user using configuration
	defaultUser := models.NewDefaultUser()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up routes
	mux.Handle("/authorize", &handlers.AuthorizeHandler{Store: memoryStore})
	mux.Handle("/token", handlers.NewTokenHandler(memoryStore))
	mux.Handle("/userinfo", &handlers.UserInfoHandler{Store: memoryStore})
	mux.Handle("/config", handlers.NewConfigHandler(memoryStore, defaultUser))

	// Start the server with the custom ServeMux
	startServer(serverPort, mux)
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
