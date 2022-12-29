package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// ErrorMessage displays error as json
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendError sends error in JSON format
func SendError(ctx *fiber.Ctx, err error) error {
	var e *fiber.Error
	errors.As(err, &e)

	// Send as 500 if not fiber.Error
	if e == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorMessage{
			Code:    fiber.ErrInternalServerError.Code,
			Message: err.Error(),
		})
	}

	// Send error
	return ctx.Status(e.Code).JSON(ErrorMessage{
		Code:    e.Code,
		Message: e.Message,
	})
}
