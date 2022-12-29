package handler

import (
	"backend/claims"
	"backend/errors"
	"backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// FindSessionHandler finds session based on parameters.
func FindSessionHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.FindSessionBody)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, count, err := services.FindSession(body, session)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	// Return json
	return ctx.JSON(fiber.Map{
		"count":    count,
		"offset":   body.Offset,
		"sessions": res,
	})
}

// RevokeSessionHandler revokes sessions from refreshing JWT-token.
func RevokeSessionHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.RevokeSessionBody)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	// Revoke session
	res, _, err := services.RevokeSession(body, session)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	// Return json
	return ctx.JSON(res)
}
