package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
	once       sync.Once
)

// InitKeys initializes the RSA key pair for JWT signing
func InitKeys() error {
	var err error
	once.Do(func() {
		// Generate RSA key pair
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return
		}
		publicKey = &privateKey.PublicKey
		keyID = "mock-key-1"
	})
	return err
}

// GenerateIDToken creates a signed JWT ID token
func GenerateIDToken(issuer, clientID, sub string) (string, error) {
	if privateKey == nil {
		if err := InitKeys(); err != nil {
			return "", err
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   issuer,
		"sub":   sub,
		"aud":   clientID,
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
		"nonce": generateNonce(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	return token.SignedString(privateKey)
}

// GenerateAccessToken creates a signed JWT access token
func GenerateAccessToken(issuer, clientID, sub string, scopes []string) (string, error) {
	if privateKey == nil {
		if err := InitKeys(); err != nil {
			return "", err
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   issuer,
		"sub":   sub,
		"aud":   clientID,
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
		"scope": scopes,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	return token.SignedString(privateKey)
}

// GetJWKS returns the JSON Web Key Set
func GetJWKS() (map[string]interface{}, error) {
	if publicKey == nil {
		if err := InitKeys(); err != nil {
			return nil, err
		}
	}

	// Encode the public key components
	nBytes := publicKey.N.Bytes()
	eBytes := big.NewInt(int64(publicKey.E)).Bytes()

	n := base64.RawURLEncoding.EncodeToString(nBytes)
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	jwk := map[string]interface{}{
		"kty": "RSA",
		"use": "sig",
		"kid": keyID,
		"alg": "RS256",
		"n":   n,
		"e":   e,
	}

	jwks := map[string]interface{}{
		"keys": []interface{}{jwk},
	}

	return jwks, nil
}

// VerifyToken verifies a JWT token and returns the claims
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	if publicKey == nil {
		if err := InitKeys(); err != nil {
			return nil, err
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateNonce generates a random nonce for the token
func generateNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// This should never happen with crypto/rand, but we handle it for safety
		panic("failed to generate random bytes: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// GetPublicKey returns the public key (for testing purposes)
func GetPublicKey() (*rsa.PublicKey, error) {
	if publicKey == nil {
		if err := InitKeys(); err != nil {
			return nil, err
		}
	}
	return publicKey, nil
}

// GetPublicKeyPEM returns the public key in PEM format
func GetPublicKeyPEM() (string, error) {
	if publicKey == nil {
		if err := InitKeys(); err != nil {
			return "", err
		}
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil
}

// MarshalJWKS returns the JWKS as a JSON byte array
func MarshalJWKS() ([]byte, error) {
	jwks, err := GetJWKS()
	if err != nil {
		return nil, err
	}
	return json.Marshal(jwks)
}
