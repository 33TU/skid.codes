package handler

import (
	"backend/claims"
	"backend/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ip2location/ip2location-go/v9"
)

// LoginHandler logins and creates session into db returning JWT token.
func LoginHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.LoginBody)
	lookup := ctx.Locals("lookup").(*ip2location.IP2Locationrecord)

	// Login
	res, auth, cookie, err := services.Login(body, lookup, ctx.IP())
	if err != nil {
		return err
	}

	// Set cookie
	ctx.Cookie(cookie)

	// Return response
	return ctx.JSON(fiber.Map{
		"res":  res,
		"auth": auth,
	})
}

// RefreshHandler refreshes the JWT token.
func RefreshHandler(ctx *fiber.Ctx) error {
	session := ctx.Locals("refresh").(*jwt.Token).Claims.(*claims.RefreshClaims)

	// Refresh
	res, auth, cookie, err := services.Refresh(session, ctx.IP())
	if err != nil {
		return err
	}

	// Set cookie
	ctx.Cookie(cookie)

	// Return response
	return ctx.JSON(fiber.Map{
		"res":  res,
		"auth": auth,
	})
}

// LogoutHandler revokes session.
func LogoutHandler(ctx *fiber.Ctx) error {
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	// Revoke session
	res, err := services.RevokeSession(&services.RevokeSessionBody{
		SessionID: session.SessionID,
	}, session)

	if err != nil {
		return err
	}

	// Send expired cookie for refresh_token
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "expired",
		Expires:  time.Time{},
		SameSite: "Strict",
		Path:     "/api/auth/refresh",
		HTTPOnly: true,
	})

	return ctx.JSON(res)
}
