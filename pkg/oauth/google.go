package oauth

import (
	"net/url"
	"time"

	"github.com/chrisw-dev/golang-mock-oauth2-server/internal/store"
)

// GoogleProvider implements the Provider interface for Google OAuth2
type GoogleProvider struct {
	Store *store.MemoryStore
}

// NewGoogleProvider creates a new Google OAuth2 provider instance
func NewGoogleProvider(store *store.MemoryStore) *GoogleProvider {
	return &GoogleProvider{Store: store}
}

// GenerateAuthURL creates an authorization URL for the OAuth2 flow
func (p *GoogleProvider) GenerateAuthURL(clientID, redirectURI, scope, state string) string {
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", scope)
	params.Set("response_type", "code")
	if state != "" {
		params.Set("state", state)
	}

	return "/authorize?" + params.Encode()
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (p *GoogleProvider) ExchangeCodeForToken(code string) (map[string]interface{}, error) {
	// Simulate token exchange
	authRequest, exists := p.Store.GetAuthCode(code)
	if !exists {
		return nil, &Error{Code: "invalid_grant", Description: "Invalid authorization code"}
	}

	if authRequest.Expiration.Before(time.Now()) {
		return nil, &Error{Code: "invalid_grant", Description: "Authorization code expired"}
	}

	token := map[string]interface{}{
		"access_token":  "mock-access-token",
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": "mock-refresh-token",
		"id_token":      "mock-id-token",
	}

	return token, nil
}

// GetUserInfo retrieves user information using the provided access token
func (p *GoogleProvider) GetUserInfo(accessToken string) (map[string]interface{}, error) {
	// Simulate user info retrieval
	userInfo, exists := p.Store.GetUserInfoByToken(accessToken)
	if !exists {
		return nil, &Error{Code: "invalid_token", Description: "Invalid access token"}
	}

	return map[string]interface{}{
		"sub":            userInfo.Sub,
		"name":           userInfo.Name,
		"given_name":     userInfo.GivenName,
		"family_name":    userInfo.FamilyName,
		"email":          userInfo.Email,
		"email_verified": userInfo.EmailVerified,
		"picture":        userInfo.Picture,
		"locale":         userInfo.Locale,
		"hd":             userInfo.HD,
	}, nil
}

// Error represents an OAuth2 error with an error code and description
type Error struct {
	Code        string
	Description string
}

func (e *Error) Error() string {
	return e.Code + ": " + e.Description
}
