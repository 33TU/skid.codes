package handler

import (
	"backend/claims"
	"backend/errors"
	"backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ip2location/ip2location-go/v9"
)

// LoginHandler logins and creates session into db returning JWT token.
func LoginHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.LoginRequest)
	lookup := ctx.Locals("lookup").(*ip2location.IP2Locationrecord)

	// Login
	res, auth, cookie, err := services.Login(body, lookup, ctx.IP())
	if err != nil {
		return errors.SendError(ctx, err)
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
		return errors.SendError(ctx, err)
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
	res, cookie, err := services.RevokeSession(&services.RevokeSessionRequest{
		SessionID: session.SessionID,
	}, session)

	if err != nil {
		return errors.SendError(ctx, err)
	}

	// Set cookie
	ctx.Cookie(cookie)

	return ctx.JSON(res)
}
