package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("MOCK_OAUTH_PORT", "9090")
	os.Setenv("MOCK_USER_EMAIL", "test@example.com")
	os.Setenv("MOCK_USER_NAME", "Test User")
	os.Setenv("MOCK_TOKEN_EXPIRY", "7200")
	defer os.Unsetenv("MOCK_OAUTH_PORT")
	defer os.Unsetenv("MOCK_USER_EMAIL")
	defer os.Unsetenv("MOCK_USER_NAME")
	defer os.Unsetenv("MOCK_TOKEN_EXPIRY")

	config := LoadConfig()

	if config.Port != 9090 {
		t.Errorf("expected Port to be 9090, got %d", config.Port)
	}
	if config.MockUserEmail != "test@example.com" {
		t.Errorf("expected MockUserEmail to be 'test@example.com', got '%s'", config.MockUserEmail)
	}
	if config.MockUserName != "Test User" {
		t.Errorf("expected MockUserName to be 'Test User', got '%s'", config.MockUserName)
	}
	if config.MockTokenExpiry != 7200 {
		t.Errorf("expected MockTokenExpiry to be 7200, got %d", config.MockTokenExpiry)
	}
}

func TestUpdateConfig(t *testing.T) {
	config := LoadConfig()

	newConfig := map[string]interface{}{
		"port":              8081,
		"mock_user_email":   "updated@example.com",
		"mock_user_name":    "Updated User",
		"mock_token_expiry": 3600,
	}

	config.UpdateConfig(newConfig)

	if config.Port != 8081 {
		t.Errorf("expected Port to be 8081, got %d", config.Port)
	}
	if config.MockUserEmail != "updated@example.com" {
		t.Errorf("expected MockUserEmail to be 'updated@example.com', got '%s'", config.MockUserEmail)
	}
	if config.MockUserName != "Updated User" {
		t.Errorf("expected MockUserName to be 'Updated User', got '%s'", config.MockUserName)
	}
	if config.MockTokenExpiry != 3600 {
		t.Errorf("expected MockTokenExpiry to be 3600, got %d", config.MockTokenExpiry)
	}
}

func TestGetConfig(t *testing.T) {
	config := LoadConfig()

	newConfig := map[string]interface{}{
		"port":              8082,
		"mock_user_email":   "getconfig@example.com",
		"mock_user_name":    "GetConfig User",
		"mock_token_expiry": 1800,
	}

	config.UpdateConfig(newConfig)
	retrievedConfig := config.GetConfig()

	if retrievedConfig.Port != 8082 ||
		retrievedConfig.MockUserEmail != "getconfig@example.com" ||
		retrievedConfig.MockUserName != "GetConfig User" ||
		retrievedConfig.MockTokenExpiry != 1800 {

		t.Errorf("expected retrievedConfig to match updated config, got Port: %d, MockUserEmail: %s, MockUserName: %s, MockTokenExpiry: %d", retrievedConfig.Port, retrievedConfig.MockUserEmail, retrievedConfig.MockUserName, retrievedConfig.MockTokenExpiry)
	}
}
