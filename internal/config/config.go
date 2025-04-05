package config

import (
	"os"
	"strconv"
	"sync"
)

type ServerConfig struct {
	Port            int
	MockUserEmail   string
	MockUserName    string
	MockTokenExpiry int
	mu              sync.RWMutex
}

var defaultConfig = ServerConfig{
	Port:            8080,
	MockUserEmail:   "testuser@example.com",
	MockUserName:    "Test User",
	MockTokenExpiry: 3600,
}

// Updated to avoid copying the lock value by using a pointer to `ServerConfig`.
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

	return config
}

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
}

// Updated `GetConfig` to return a new struct without the mutex field.
func (c *ServerConfig) GetConfig() ServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return ServerConfig{
		Port:            c.Port,
		MockUserEmail:   c.MockUserEmail,
		MockUserName:    c.MockUserName,
		MockTokenExpiry: c.MockTokenExpiry,
	}
}
