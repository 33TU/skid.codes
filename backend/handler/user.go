package handler

import (
	"backend/claims"
	"backend/errors"
	"backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// GetUserHandler gets user based on parameters.
func GetUserHandler(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	res, err := services.GetUser(username)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}

// FindPasteHandler finds user based on parameters.
func FindUserHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.FindUserRequest)

	res, count, err := services.FindUser(body)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	// Return json
	return ctx.JSON(fiber.Map{
		"count":  count,
		"offset": body.Offset,
		"users":  res,
	})
}

// UpdatePasteHandler updates user's details.
func UpdateUserHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.UpdateUserRequest)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, err := services.UpdateUser(body, session)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}

// CreatePasteHandler creates new user.
func CreateUserHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.CreateUserRequest)

	res, err := services.CreateUser(body)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}
