package middlewares

import (
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/validator"

	"github.com/gofiber/fiber/v2"
)

type defaultSetter interface {
	SetDefaults()
}

func validate[T any](c *fiber.Ctx, source string, parser func(dest any) error, localKey string) error {
	var data T

	if err := parser(&data); err != nil {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			"INvalid"+source,
			response.ErrBadRequest,
			err,
			nil,
		)
	}

	if def, ok := any(&data).(defaultSetter); ok {
		def.SetDefaults()
	}

	if detail, err := validator.ValidateStruct(data); err != nil {
		var errorDetail any = err
		if detail != nil {
			errorDetail = detail
		}

		return response.Error(
			c,
			fiber.StatusBadRequest,
			source+"validation error",
			response.ErrValidation,
			errorDetail,
			nil,
		)
	}

	c.Locals(localKey, data)

	return c.Next()
}

func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "body", c.BodyParser, "body")
	}
}

func ValidateQuery[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "query", c.QueryParser, "query")
	}
}

func ValidateParams[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "params", c.ParamsParser, "params")
	}
}
