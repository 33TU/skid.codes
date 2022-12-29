package middleware

import (
	"backend/claims"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v4"
)

// limitReached returns fiber error for limit reached
func limitReached(ctx *fiber.Ctx) error {
	return fiber.ErrTooManyRequests
}

// AuthLimiter limits by sender email.
func AuthLimiter(max int, exp time.Duration, skipFailed bool) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                max,
		Expiration:         exp,
		LimitReached:       limitReached,
		SkipFailedRequests: skipFailed,

		KeyGenerator: func(ctx *fiber.Ctx) string {
			session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)
			return session.Email
		},
	})
}

// RefreshLimiter limits by sender id.
func RefreshLimiter(max int, exp time.Duration, skipFailed bool) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                max,
		Expiration:         exp,
		LimitReached:       limitReached,
		SkipFailedRequests: skipFailed,

		KeyGenerator: func(ctx *fiber.Ctx) string {
			session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.RefreshClaims)
			return strconv.Itoa(session.SessionID)
		},
	})
}

// IPLimiter limits by sender IP.
func IPLimiter(max int, exp time.Duration, skipFailed bool) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                max,
		Expiration:         exp,
		LimitReached:       limitReached,
		SkipFailedRequests: skipFailed,

		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
}
