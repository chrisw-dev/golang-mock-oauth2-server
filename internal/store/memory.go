package store

import (
	"log"
	"sync"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/types"
)

// Store interface defines methods for storing and retrieving OAuth2 data
type Store interface {
	// Auth code methods
	StoreAuthCode(code string, request *models.AuthRequest)
	GetAuthCode(code string) (*models.AuthRequest, bool)
	RemoveAuthCode(code string)

	// Token methods
	StoreToken(token string, clientID string)
	GetClientIDByToken(token string) (string, bool)

	// Config methods
	StoreTokenConfig(config map[string]interface{})
	GetTokenConfig() map[string]interface{}
	StoreErrorScenario(scenario types.ErrorScenario)
	GetErrorScenario(endpoint string) (*types.ErrorScenario, bool)
	ClearErrorScenario(endpoint string)
}

// MemoryStore implements Store using in-memory storage
type MemoryStore struct {
	mu            sync.RWMutex
	authCodes     map[string]*models.AuthRequest
	tokens        map[string]string // token -> clientID
	tokenConfig   map[string]interface{}
	errorScenario *types.ErrorScenario
}

// NewMemoryStore creates a new memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		authCodes:   make(map[string]*models.AuthRequest),
		tokens:      make(map[string]string),
		tokenConfig: make(map[string]interface{}),
	}
}

// StoreAuthCode stores an authorization code
func (s *MemoryStore) StoreAuthCode(code string, request *models.AuthRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authCodes[code] = request
}

// GetAuthCode retrieves an authorization code by its value
func (s *MemoryStore) GetAuthCode(code string) (*models.AuthRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	request, exists := s.authCodes[code]
	return request, exists
}

// RemoveAuthCode removes an authorization code
func (s *MemoryStore) RemoveAuthCode(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.authCodes, code)
}

// StoreToken stores a token with its associated client ID
func (s *MemoryStore) StoreToken(token string, clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = clientID
}

// GetClientIDByToken retrieves the client ID associated with a token
func (s *MemoryStore) GetClientIDByToken(token string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	clientID, exists := s.tokens[token]
	return clientID, exists
}

// StoreTokenConfig saves customized token configuration
func (s *MemoryStore) StoreTokenConfig(config map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Replace the entire config
	s.tokenConfig = config
}

// GetTokenConfig retrieves the current token configuration
func (s *MemoryStore) GetTokenConfig() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent race conditions
	if s.tokenConfig == nil {
		return nil
	}

	config := make(map[string]interface{}, len(s.tokenConfig))
	for k, v := range s.tokenConfig {
		config[k] = v
	}
	return config
}

// StoreErrorScenario stores an error scenario configuration
func (s *MemoryStore) StoreErrorScenario(scenario types.ErrorScenario) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errorScenario = &scenario
}

// GetErrorScenario retrieves error scenario for the specified endpoint
func (s *MemoryStore) GetErrorScenario(endpoint string) (*types.ErrorScenario, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.errorScenario == nil {
		return nil, false
	}

	// Only return if this scenario is for the requested endpoint
	if s.errorScenario.Endpoint == endpoint {
		return s.errorScenario, true
	}

	return nil, false
}

// ClearErrorScenario removes the error scenario for the specified endpoint
func (s *MemoryStore) ClearErrorScenario(endpoint string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.errorScenario != nil && s.errorScenario.Endpoint == endpoint {
		s.errorScenario = nil
	}
}

// GetUserInfoByToken retrieves user information based on a token
func (s *MemoryStore) GetUserInfoByToken(token string) (*models.UserInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(token) > 4 {
		log.Printf("Store: Looking up token: %s...", token[:4])
	} else {
		log.Printf("Store: Looking up token (short token)")
	}

	clientID, exists := s.tokens[token]
	if !exists {
		log.Printf("Store: Token not found in tokens map")
		return nil, false
	}
	log.Printf("Store: Found clientID for token: %s", clientID)

	// Instead of trying to look up the auth request (which is removed after token exchange),
	// we'll generate a user info object directly from the clientID

	// Extract the base clientID (removing any prefixes/suffixes added by token generation)
	baseClientID := clientID

	// Generate a default user info
	userInfo := &models.UserInfo{
		Sub:           baseClientID,
		ID:            "1234",
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Email:         baseClientID + "@example.com",
		Picture:       "https://example.com/photo.jpg",
		EmailVerified: true,
	}

	// Check if we have custom user info configuration
	if userInfoConfig, ok := s.tokenConfig["user_info"].(map[string]interface{}); ok {
		log.Printf("Store: Found custom user info configuration")

		// Override default values with configured ones
		if sub, ok := userInfoConfig["sub"].(string); ok {
			userInfo.Sub = sub
		}
		if name, ok := userInfoConfig["name"].(string); ok {
			userInfo.Name = name
		}
		if email, ok := userInfoConfig["email"].(string); ok {
			userInfo.Email = email
		}
		if verified, ok := userInfoConfig["email_verified"].(bool); ok {
			userInfo.EmailVerified = verified
		}
		if givenName, ok := userInfoConfig["given_name"].(string); ok {
			userInfo.GivenName = givenName
		}
		if familyName, ok := userInfoConfig["family_name"].(string); ok {
			userInfo.FamilyName = familyName
		}
		if picture, ok := userInfoConfig["picture"].(string); ok {
			userInfo.Picture = picture
		}
	}

	log.Printf("Store: Generated user info with email: %s", userInfo.Email)
	return userInfo, true
}
