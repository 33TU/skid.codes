package middleware

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
)

var (
	validate = validator.New()
)

// Struct validates a structs exposed fields. If invalid error message is built.
func structValidate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Fprintln(buf, err)
	}

	for _, err := range err.(validator.ValidationErrors) {
		fmt.Fprintln(buf, err.StructField(), "is not valid. Reason: tag", err.Tag(), "is not met.")
	}

	return errors.New(buf.String())
}

// Validate validates T and on success stores it to ctx's locals as "body".
func Validate[T any]() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body T

		if err := ctx.BodyParser(&body); err != nil {
			return err
		}

		if err := structValidate(body); err != nil {
			return err
		}

		ctx.Locals("body", &body)
		return ctx.Next()
	}
}
