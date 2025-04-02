package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestUserInfoMarshaling(t *testing.T) {
	// Create a sample user
	user := UserInfo{
		Sub:           "test-id-123",
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Email:         "test@example.com",
		EmailVerified: true,
		Picture:       "https://example.com/picture.jpg",
		Locale:        "en-US",
		HD:            "example.com",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal UserInfo: %v", err)
	}

	// Unmarshal back to struct
	var unmarshaled UserInfo
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserInfo: %v", err)
	}

	// Verify unmarshaled data matches original
	if !reflect.DeepEqual(user, unmarshaled) {
		t.Errorf("Unmarshaled data doesn't match original. Got: %+v, Want: %+v", unmarshaled, user)
	}
}

func TestUserInfoJSONFieldNames(t *testing.T) {
	// Create a sample user
	user := UserInfo{
		Sub:           "test-id-123",
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Email:         "test@example.com",
		EmailVerified: true,
		Picture:       "https://example.com/picture.jpg",
		Locale:        "en-US",
		HD:            "example.com",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal UserInfo: %v", err)
	}

	// Convert to string for inspection
	jsonStr := string(jsonData)

	// Check that the JSON contains the expected field names
	expectedFields := []string{
		`"sub"`,
		`"name"`,
		`"given_name"`,
		`"family_name"`,
		`"email"`,
		`"email_verified"`,
		`"picture"`,
		`"locale"`,
		`"hd"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON string doesn't contain expected field %s. JSON: %s", field, jsonStr)
		}
	}
}

func TestNewDefaultUser(t *testing.T) {
	// Get default user
	user := NewDefaultUser()

	// Check default values
	if user.Sub != "123456789" {
		t.Errorf("Expected default Sub to be '123456789', got '%s'", user.Sub)
	}
	if user.Name != "Test User" {
		t.Errorf("Expected default Name to be 'Test User', got '%s'", user.Name)
	}
	if user.Email != "testuser@example.com" {
		t.Errorf("Expected default Email to be 'testuser@example.com', got '%s'", user.Email)
	}
	if !user.EmailVerified {
		t.Error("Expected default EmailVerified to be true")
	}
}

func TestUserInfoClone(t *testing.T) {
	// Create a test user
	original := &UserInfo{
		Sub:           "test-id-123",
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Email:         "test@example.com",
		EmailVerified: true,
		Picture:       "https://example.com/picture.jpg",
		Locale:        "en-US",
		HD:            "example.com",
	}

	// Clone the user
	cloned := original.Clone()

	// Verify cloned data matches original
	if !reflect.DeepEqual(original, cloned) {
		t.Errorf("Cloned data doesn't match original. Got: %+v, Want: %+v", cloned, original)
	}

	// Verify it's a deep copy by modifying the clone and checking the original
	cloned.Name = "Modified Name"
	cloned.Email = "modified@example.com"

	if original.Name == cloned.Name {
		t.Error("Clone is not a deep copy - name was modified in original")
	}
	if original.Email == cloned.Email {
		t.Error("Clone is not a deep copy - email was modified in original")
	}

	// Test cloning nil
	var nilUser *UserInfo
	nilClone := nilUser.Clone()
	if nilClone != nil {
		t.Error("Cloning nil should return nil")
	}
}

func TestUserInfoMerge(t *testing.T) {
	testCases := []struct {
		name     string
		base     *UserInfo
		other    *UserInfo
		expected *UserInfo
	}{
		{
			name: "Merge all fields",
			base: &UserInfo{
				Sub:           "base-id",
				Name:          "Base User",
				GivenName:     "Base",
				FamilyName:    "User",
				Email:         "base@example.com",
				EmailVerified: false,
				Picture:       "https://example.com/base.jpg",
			},
			other: &UserInfo{
				Sub:           "other-id",
				Name:          "Other User",
				GivenName:     "Other",
				FamilyName:    "User2",
				Email:         "other@example.com",
				EmailVerified: true,
				Picture:       "https://example.com/other.jpg",
				Locale:        "fr-FR",
				HD:            "other.com",
			},
			expected: &UserInfo{
				Sub:           "other-id",
				Name:          "Other User",
				GivenName:     "Other",
				FamilyName:    "User2",
				Email:         "other@example.com",
				EmailVerified: true,
				Picture:       "https://example.com/other.jpg",
				Locale:        "fr-FR",
				HD:            "other.com",
			},
		},
		{
			name: "Merge some fields",
			base: &UserInfo{
				Sub:           "base-id",
				Name:          "Base User",
				Email:         "base@example.com",
				EmailVerified: false,
			},
			other: &UserInfo{
				Name: "Other User",
				// Email is empty, should not overwrite
				EmailVerified: true, // Boolean should always update
			},
			expected: &UserInfo{
				Sub:           "base-id",
				Name:          "Other User",
				Email:         "base@example.com",
				EmailVerified: true,
			},
		},
		{
			name: "Merge with nil",
			base: &UserInfo{
				Sub:  "base-id",
				Name: "Base User",
			},
			other: nil,
			expected: &UserInfo{
				Sub:  "base-id",
				Name: "Base User",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.base.Merge(tc.other)
			if !reflect.DeepEqual(tc.base, tc.expected) {
				t.Errorf("Merged result doesn't match expected. Got: %+v, Want: %+v", tc.base, tc.expected)
			}
		})
	}
}

func TestUpdateUserFromConfig(t *testing.T) {
	testCases := []struct {
		name     string
		user     *UserInfo
		config   map[string]interface{}
		expected *UserInfo
	}{
		{
			name: "Update all fields",
			user: &UserInfo{
				Sub:           "base-id",
				Name:          "Base User",
				GivenName:     "Base",
				FamilyName:    "User",
				Email:         "base@example.com",
				EmailVerified: false,
				Picture:       "https://example.com/base.jpg",
			},
			config: map[string]interface{}{
				"sub":            "config-id",
				"name":           "Config User",
				"given_name":     "Config",
				"family_name":    "User2",
				"email":          "config@example.com",
				"email_verified": true,
				"picture":        "https://example.com/config.jpg",
				"locale":         "es-ES",
				"hd":             "config.com",
			},
			expected: &UserInfo{
				Sub:           "config-id",
				Name:          "Config User",
				GivenName:     "Config",
				FamilyName:    "User2",
				Email:         "config@example.com",
				EmailVerified: true,
				Picture:       "https://example.com/config.jpg",
				Locale:        "es-ES",
				HD:            "config.com",
			},
		},
		{
			name: "Update some fields",
			user: &UserInfo{
				Sub:           "base-id",
				Name:          "Base User",
				Email:         "base@example.com",
				EmailVerified: false,
			},
			config: map[string]interface{}{
				"name":           "Config User",
				"email_verified": true,
			},
			expected: &UserInfo{
				Sub:           "base-id",
				Name:          "Config User",
				Email:         "base@example.com",
				EmailVerified: true,
			},
		},
		{
			name: "Update with nil config",
			user: &UserInfo{
				Sub:  "base-id",
				Name: "Base User",
			},
			config: nil,
			expected: &UserInfo{
				Sub:  "base-id",
				Name: "Base User",
			},
		},
		{
			name: "Update with wrong types",
			user: &UserInfo{
				Sub:  "base-id",
				Name: "Base User",
			},
			config: map[string]interface{}{
				"sub":            123, // Not a string, should be ignored
				"name":           "Config User",
				"email_verified": "true", // Not a bool, should be ignored
			},
			expected: &UserInfo{
				Sub:  "base-id",
				Name: "Config User",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			UpdateUserFromConfig(tc.user, tc.config)
			if !reflect.DeepEqual(tc.user, tc.expected) {
				t.Errorf("Updated user doesn't match expected. Got: %+v, Want: %+v", tc.user, tc.expected)
			}
		})
	}
}

// Add a test for omitempty behavior
func TestUserInfoOmitEmptyFields(t *testing.T) {
	// Create a user with only required fields
	user := &UserInfo{
		Sub:           "test-id",
		Name:          "Test User",
		Email:         "test@example.com",
		EmailVerified: true,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal UserInfo: %v", err)
	}

	jsonStr := string(jsonData)

	// Check that optional empty fields are omitted
	if strings.Contains(jsonStr, `"locale":`) {
		t.Error("Empty locale field should be omitted in JSON")
	}
	if strings.Contains(jsonStr, `"hd":`) {
		t.Error("Empty hd field should be omitted in JSON")
	}
}
