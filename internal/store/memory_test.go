package store

import (
	"reflect"
	"testing"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/types"
)

func TestMemoryStore_AuthCodeMethods(t *testing.T) {
	store := NewMemoryStore()
	authCode := "test-code"
	authRequest := &models.AuthRequest{ClientID: "test-client"}

	// Test StoreAuthCode
	store.StoreAuthCode(authCode, authRequest)
	storedRequest, exists := store.GetAuthCode(authCode)
	if !exists || !reflect.DeepEqual(storedRequest, authRequest) {
		t.Errorf("expected stored auth request to match, got %+v", storedRequest)
	}

	// Test RemoveAuthCode
	store.RemoveAuthCode(authCode)
	_, exists = store.GetAuthCode(authCode)
	if exists {
		t.Errorf("expected auth code to be removed")
	}
}

func TestMemoryStore_TokenMethods(t *testing.T) {
	store := NewMemoryStore()
	token := "test-token"
	clientID := "test-client"

	// Test StoreToken
	store.StoreToken(token, clientID)
	storedClientID, exists := store.GetClientIDByToken(token)
	if !exists || storedClientID != clientID {
		t.Errorf("expected stored client ID to match, got %s", storedClientID)
	}
}

func TestMemoryStore_ConfigMethods(t *testing.T) {
	store := NewMemoryStore()
	tokenConfig := map[string]interface{}{
		"key": "value",
	}

	// Test StoreTokenConfig
	store.StoreTokenConfig(tokenConfig)
	storedConfig := store.GetTokenConfig()
	if !reflect.DeepEqual(storedConfig, tokenConfig) {
		t.Errorf("expected stored config to match, got %+v", storedConfig)
	}
}

func TestMemoryStore_ErrorScenarioMethods(t *testing.T) {
	store := NewMemoryStore()
	scenario := types.ErrorScenario{
		Endpoint:    "test-endpoint",
		StatusCode:  400,
		ErrorCode:   "invalid_request",
		Description: "Invalid request",
	}

	// Test StoreErrorScenario
	store.StoreErrorScenario(scenario)
	storedScenario, exists := store.GetErrorScenario("test-endpoint")
	if !exists || !reflect.DeepEqual(storedScenario, &scenario) {
		t.Errorf("expected stored scenario to match, got %+v", storedScenario)
	}

	// Test ClearErrorScenario
	store.ClearErrorScenario("test-endpoint")
	_, exists = store.GetErrorScenario("test-endpoint")
	if exists {
		t.Errorf("expected error scenario to be cleared")
	}
}

func TestMemoryStore_GetUserInfoByToken(t *testing.T) {
	store := NewMemoryStore()
	token := "test-token"
	clientID := "test-client"
	authRequest := &models.AuthRequest{ClientID: clientID}

	store.StoreToken(token, clientID)
	store.StoreAuthCode(clientID, authRequest)

	userInfo, exists := store.GetUserInfoByToken(token)
	if !exists || userInfo.Sub != clientID || userInfo.Email != clientID+"@example.com" {
		t.Errorf("expected user info to match, got %+v", userInfo)
	}
}
