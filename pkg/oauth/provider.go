package oauth

// Provider defines the interface for OAuth2 providers
type Provider interface {
	GenerateAuthURL(clientID, redirectURI, scope, state string) string
	ExchangeCodeForToken(code string) (map[string]interface{}, error)
	GetUserInfo(accessToken string) (map[string]interface{}, error)
}
