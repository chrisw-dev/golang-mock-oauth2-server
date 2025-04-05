package store

import (
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

// Reverted `GetAuthCode` to return both `*models.AuthRequest` and a boolean.
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

	clientID, exists := s.tokens[token]
	if !exists {
		return nil, false
	}

	authRequest, exists := s.authCodes[clientID]
	if !exists {
		return nil, false
	}

	// Transform AuthRequest into UserInfo
	userInfo := &models.UserInfo{
		Sub:           clientID,                              // Using clientID as a unique identifier
		Name:          "Generated User",                      // Placeholder name
		Email:         authRequest.ClientID + "@example.com", // Placeholder email
		EmailVerified: true,                                  // Default to verified
	}

	return userInfo, true
}
