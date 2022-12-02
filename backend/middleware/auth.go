package middleware

import (
	"backend/claims"
	"backend/secret"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// Auth middleware for auth_token
func Auth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningMethod: secret.SigningMethod(),
		SigningKey:    secret.PublicKey(),
		ContextKey:    "auth",
		Claims:        &claims.AuthClaims{},
	})
}

// Refresh middleware for refresh_token
func Refresh() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningMethod: secret.SigningMethod(),
		SigningKey:    secret.PublicKey(),
		ContextKey:    "refresh",
		Claims:        &claims.RefreshClaims{},
		TokenLookup:   "cookie:refresh_token",
	})
}
