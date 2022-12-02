package claims

import (
	"backend/secret"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	authTokenTimeout    = time.Minute * 5
	refreshTokenTimeout = time.Hour * 72
)

type AuthClaims struct {
	SessionID int    `json:"sid"`
	UserID    int    `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	SessionID int `json:"sid"`
	jwt.StandardClaims
}

// Create refresh token for this auth token
func (ac *AuthClaims) Refresh() *RefreshClaims {
	return &RefreshClaims{
		SessionID: ac.SessionID,
	}
}

// CreateToken generate token strings
func CreateTokens(auth *AuthClaims, refresh *RefreshClaims) (string, string, time.Time, error) {
	// Set exp
	auth.ExpiresAt = time.Now().Add(authTokenTimeout).Unix()
	refreshExpires := time.Now().Add(refreshTokenTimeout)
	refresh.ExpiresAt = refreshExpires.Unix()

	// Create tokens
	authToken := jwt.NewWithClaims(jwt.SigningMethodRS256, auth)
	refeshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refresh)

	// Generate encoded tokens and send them as response.
	authTokenStr, err := authToken.SignedString(secret.PrivateKey())
	if err != nil {
		return "", "", refreshExpires, err
	}

	refeshTokenStr, err := refeshToken.SignedString(secret.PrivateKey())
	if err != nil {
		return "", "", refreshExpires, err
	}

	return authTokenStr, refeshTokenStr, refreshExpires, nil
}
