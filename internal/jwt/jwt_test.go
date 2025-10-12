package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestInitKeys(t *testing.T) {
	err := InitKeys()
	if err != nil {
		t.Fatalf("Failed to initialize keys: %v", err)
	}

	if privateKey == nil {
		t.Error("Private key should not be nil after initialization")
	}

	if publicKey == nil {
		t.Error("Public key should not be nil after initialization")
	}

	if keyID == "" {
		t.Error("Key ID should not be empty after initialization")
	}
}

func TestGenerateIDToken(t *testing.T) {
	err := InitKeys()
	if err != nil {
		t.Fatalf("Failed to initialize keys: %v", err)
	}

	issuer := "http://localhost:8080"
	clientID := "test-client"
	sub := "user-123"
	email := "user-123@example.com"

	tokenString, err := GenerateIDToken(issuer, clientID, sub, email)
	if err != nil {
		t.Fatalf("Failed to generate ID token: %v", err)
	}

	if tokenString == "" {
		t.Error("Token string should not be empty")
	}

	// Verify the token can be parsed
	claims, err := VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	// Check claims
	if claims["iss"] != issuer {
		t.Errorf("Expected issuer %s, got %v", issuer, claims["iss"])
	}

	if claims["sub"] != sub {
		t.Errorf("Expected subject %s, got %v", sub, claims["sub"])
	}

	if claims["aud"] != clientID {
		t.Errorf("Expected audience %s, got %v", clientID, claims["aud"])
	}

	if claims["email"] != email {
		t.Errorf("Expected email %s, got %v", email, claims["email"])
	}
}

func TestGenerateAccessToken(t *testing.T) {
	err := InitKeys()
	if err != nil {
		t.Fatalf("Failed to initialize keys: %v", err)
	}

	issuer := "http://localhost:8080"
	clientID := "test-client"
	sub := "user-123"
	scopes := []string{"openid", "email", "profile"}

	tokenString, err := GenerateAccessToken(issuer, clientID, sub, scopes)
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	if tokenString == "" {
		t.Error("Token string should not be empty")
	}

	// Verify the token can be parsed
	claims, err := VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	// Check claims
	if claims["iss"] != issuer {
		t.Errorf("Expected issuer %s, got %v", issuer, claims["iss"])
	}

	if claims["sub"] != sub {
		t.Errorf("Expected subject %s, got %v", sub, claims["sub"])
	}
}

func TestGetJWKS(t *testing.T) {
	err := InitKeys()
	if err != nil {
		t.Fatalf("Failed to initialize keys: %v", err)
	}

	jwks, err := GetJWKS()
	if err != nil {
		t.Fatalf("Failed to get JWKS: %v", err)
	}

	if jwks == nil {
		t.Error("JWKS should not be nil")
	}

	keys, ok := jwks["keys"].([]interface{})
	if !ok {
		t.Fatal("JWKS should have a 'keys' array")
	}

	if len(keys) == 0 {
		t.Error("JWKS keys array should not be empty")
	}

	// Check the first key
	key, ok := keys[0].(map[string]interface{})
	if !ok {
		t.Fatal("Key should be a map")
	}

	if key["kty"] != "RSA" {
		t.Errorf("Expected kty to be RSA, got %v", key["kty"])
	}

	if key["use"] != "sig" {
		t.Errorf("Expected use to be sig, got %v", key["use"])
	}

	if key["alg"] != "RS256" {
		t.Errorf("Expected alg to be RS256, got %v", key["alg"])
	}
}

func TestVerifyToken(t *testing.T) {
	err := InitKeys()
	if err != nil {
		t.Fatalf("Failed to initialize keys: %v", err)
	}

	// Generate a valid token
	issuer := "http://localhost:8080"
	clientID := "test-client"
	sub := "user-123"
	email := "user-123@example.com"

	tokenString, err := GenerateIDToken(issuer, clientID, sub, email)
	if err != nil {
		t.Fatalf("Failed to generate ID token: %v", err)
	}

	// Verify the token
	claims, err := VerifyToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if claims == nil {
		t.Error("Claims should not be nil")
	}

	// Test with invalid token
	invalidToken := "invalid.token.string"
	_, err = VerifyToken(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestTokenFormat(t *testing.T) {
	err := InitKeys()
	if err != nil {
		t.Fatalf("Failed to initialize keys: %v", err)
	}

	issuer := "http://localhost:8080"
	clientID := "test-client"
	sub := "user-123"
	email := "user-123@example.com"

	tokenString, err := GenerateIDToken(issuer, clientID, sub, email)
	if err != nil {
		t.Fatalf("Failed to generate ID token: %v", err)
	}

	// Parse token to check it has 3 parts (header.payload.signature)
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Check that the token has the kid header
	if kid, ok := token.Header["kid"].(string); !ok || kid == "" {
		t.Error("Token should have a kid header")
	}

	// Check that the algorithm is RS256
	if alg, ok := token.Header["alg"].(string); !ok || alg != "RS256" {
		t.Errorf("Expected algorithm RS256, got %v", token.Header["alg"])
	}
}
