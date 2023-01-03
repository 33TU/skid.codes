package handler

import (
	"backend/claims"
	"backend/errors"
	"backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// FetchUserPasteHandler fetches paste based on parameters.
func FetchPasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.FetchPasteRequest)

	res, err := services.FetchPaste(body, nil)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}

// FetchUserPasteHandler fetches paste based on parameters. Passes session user id as caller.
func FetchUserPasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.FetchPasteRequest)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, err := services.FetchPaste(body, &session.UserID)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}

// FindPasteHandler finds paste based on parameters.
func FindPasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.FindPasteRequest)

	res, count, err := services.FindPaste(body, nil)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"count":  count,
		"offset": body.Offset,
		"pastes": res,
	})
}

// FindUserPasteHandler finds paste based on parameters. Passes session user id as caller.
func FindUserPasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.FindPasteRequest)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, count, err := services.FindPaste(body, &session.UserID)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"count":  count,
		"offset": body.Offset,
		"pastes": res,
	})
}

// UpdatePasteHandler updates paste's details.
func UpdatePasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.UpdatePasteRequest)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, err := services.UpdatePaste(body, session)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}

// DeletePasteHandler deletes paste.
func DeletePasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.DeletePasteRequest)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, err := services.DeletePaste(body, session)
	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}

// CreatePasteHandler creates new paste.
func CreatePasteHandler(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(*services.CreatePasteRequest)
	session := ctx.Locals("auth").(*jwt.Token).Claims.(*claims.AuthClaims)

	res, err := services.CreatePaste(body, session)

	if err != nil {
		return errors.SendError(ctx, err)
	}

	return ctx.JSON(res)
}
