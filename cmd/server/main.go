package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/config"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/handlers"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/version"
)

func main() {
	// Define command-line flags
	var port int
	var host string
	flag.IntVar(&port, "port", 0, "Port to run the server on (default: uses MOCK_OAUTH_PORT env var or 8080)")
	flag.StringVar(&host, "host", "", "Host for public URLs (default: http://localhost:[port])")
	flag.Parse()

	// Log version info on startup
	versionInfo := version.GetVersion()
	log.Printf("Starting Mock OAuth2 Server %s (commit: %s, build: %s)",
		versionInfo.Version, versionInfo.Commit, versionInfo.BuildDate)

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

	// Determine the base URL for OpenID Connect configuration
	baseURL := host
	
	// If host flag is not provided, check for IssuerURL in config (from MOCK_ISSUER_URL env var)
	if baseURL == "" {
		baseURL = cfg.IssuerURL
	}
	
	// If neither host flag nor MOCK_ISSUER_URL env var is provided, use localhost
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%d", serverPort)
	}

	log.Printf("Using issuer URL: %s", baseURL)

	// Initialize in-memory store with configuration
	memoryStore := store.NewMemoryStore()

	// Set up default user using configuration
	defaultUser := models.NewDefaultUser()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up routes
	mux.Handle("/authorize", &handlers.AuthorizeHandler{Store: memoryStore})
	mux.Handle("/token", handlers.NewTokenHandlerWithIssuer(memoryStore, baseURL))
	mux.Handle("/userinfo", &handlers.UserInfoHandler{Store: memoryStore})
	mux.Handle("/config", handlers.NewConfigHandler(memoryStore, defaultUser))
	mux.Handle("/version", handlers.NewVersionHandler())
	
	// Add OpenID Connect Discovery endpoint
	mux.Handle("/.well-known/openid-configuration", handlers.NewOpenIDConfigHandler(baseURL))
	
	// Add JWKS endpoint
	mux.Handle("/jwks", handlers.NewJWKSHandler())

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
