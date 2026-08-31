package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

// ServerConfig holds configuration parameters for the OAuth2 server
type ServerConfig struct {
	Port              int
	MockUserEmail     string
	MockUserName      string
	MockTokenExpiry   int
	IssuerURL         string
	GISEnabled        bool
	GISAllowedOrigins []string
	mu                sync.RWMutex
}

var defaultConfig = ServerConfig{
	Port:            8080,
	MockUserEmail:   "testuser@example.com",
	MockUserName:    "Test User",
	MockTokenExpiry: 3600,
	IssuerURL:       "", // Will be auto-generated if not specified
	GISEnabled:      false,
}

// LoadConfig loads server configuration from environment variables or returns defaults
func LoadConfig() *ServerConfig {
	config := &defaultConfig

	if port, exists := os.LookupEnv("MOCK_OAUTH_PORT"); exists {
		if parsedPort, err := strconv.Atoi(port); err == nil {
			config.Port = parsedPort
		}
	}

	if email, exists := os.LookupEnv("MOCK_USER_EMAIL"); exists {
		config.MockUserEmail = email
	}

	if name, exists := os.LookupEnv("MOCK_USER_NAME"); exists {
		config.MockUserName = name
	}

	if expiry, exists := os.LookupEnv("MOCK_TOKEN_EXPIRY"); exists {
		if parsedExpiry, err := strconv.Atoi(expiry); err == nil {
			config.MockTokenExpiry = parsedExpiry
		}
	}

	// Load issuer URL from environment variable
	if issuerURL, exists := os.LookupEnv("MOCK_ISSUER_URL"); exists {
		config.IssuerURL = issuerURL
	}

	// GIS compatibility is opt-in and defaults to false on any parse error
	if gisEnabled, exists := os.LookupEnv("MOCK_GIS_ENABLED"); exists {
		if parsedEnabled, err := strconv.ParseBool(gisEnabled); err == nil {
			config.GISEnabled = parsedEnabled
		}
	}

	if allowedOrigins, exists := os.LookupEnv("MOCK_GIS_ALLOWED_ORIGINS"); exists {
		config.GISAllowedOrigins = parseAllowedOrigins(allowedOrigins)
	}

	return config
}

// parseAllowedOrigins splits a comma-separated origin list, trimming whitespace and empty entries
func parseAllowedOrigins(raw string) []string {
	var origins []string
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

// UpdateConfig updates the server configuration with values from the provided map
func (c *ServerConfig) UpdateConfig(newConfig map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if port, ok := newConfig["port"].(int); ok {
		c.Port = port
	}
	if email, ok := newConfig["mock_user_email"].(string); ok {
		c.MockUserEmail = email
	}
	if name, ok := newConfig["mock_user_name"].(string); ok {
		c.MockUserName = name
	}
	if expiry, ok := newConfig["mock_token_expiry"].(int); ok {
		c.MockTokenExpiry = expiry
	}
	if issuerURL, ok := newConfig["issuer_url"].(string); ok {
		c.IssuerURL = issuerURL
	}
}

// GetConfig returns a copy of the current server configuration without the mutex
func (c *ServerConfig) GetConfig() ServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return ServerConfig{
		Port:              c.Port,
		MockUserEmail:     c.MockUserEmail,
		MockUserName:      c.MockUserName,
		MockTokenExpiry:   c.MockTokenExpiry,
		IssuerURL:         c.IssuerURL,
		GISEnabled:        c.GISEnabled,
		GISAllowedOrigins: c.GISAllowedOrigins,
	}
}
