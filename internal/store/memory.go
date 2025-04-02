package store

import (
	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/models"
)

// Store interface defines methods for storing and retrieving OAuth2 data
type Store interface {
	StoreAuthCode(code string, request *models.AuthRequest)
	GetAuthCode(code string) (*models.AuthRequest, bool)
	RemoveAuthCode(code string)
	StoreToken(token string, clientID string)
	// other methods...
}

// MemoryStore implements Store using in-memory storage
type MemoryStore struct {
	authCodes map[string]*models.AuthRequest
	tokens    map[string]string // token -> clientID
}

// NewMemoryStore creates a new memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		authCodes: make(map[string]*models.AuthRequest),
		tokens:    make(map[string]string),
	}
}

// StoreAuthCode stores an authorization code
func (s *MemoryStore) StoreAuthCode(code string, request *models.AuthRequest) {
	s.authCodes[code] = request
}

// GetAuthCode retrieves an authorization code
func (s *MemoryStore) GetAuthCode(code string) (*models.AuthRequest, bool) {
	request, exists := s.authCodes[code]
	return request, exists
}

// RemoveAuthCode removes an authorization code
func (s *MemoryStore) RemoveAuthCode(code string) {
	delete(s.authCodes, code)
}

// StoreToken stores a token with its associated client ID
func (s *MemoryStore) StoreToken(token string, clientID string) {
	s.tokens[token] = clientID
}

// Other methods as needed...
