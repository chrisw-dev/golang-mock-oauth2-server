package models

// TokenResponse represents an OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// AuthRequest represents an authorization request
type AuthRequest struct {
	ClientID    string
	RedirectURI string
	// Other fields as needed...
}
