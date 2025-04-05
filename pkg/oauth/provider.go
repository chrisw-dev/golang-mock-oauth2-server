package oauth

type OAuthProvider interface {
	GenerateAuthURL(clientID, redirectURI, scope, state string) string
	ExchangeCodeForToken(code string) (map[string]interface{}, error)
	GetUserInfo(accessToken string) (map[string]interface{}, error)
}
