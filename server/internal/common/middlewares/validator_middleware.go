package middlewares

import (
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/validator"

	"github.com/gofiber/fiber/v2"
)

type DefaultSetter interface {
	SetDefaults()
}

type SimpleValidator interface {
	Validate() error
}

type Parser interface {
	ValidateAndParse() (any, error)
}

type ParamParser interface {
	ValidateAndParseWithParams(params any) (any, error)
}

func validate[T any](
	c *fiber.Ctx,
	source string,
	parser func(dest any) error,
	localKey string,
	storeRaw bool,
	params ...any,
) error {
	var data T

	if err := parser(&data); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid "+source, response.ErrBadRequest, err)
	}

	if def, ok := any(&data).(DefaultSetter); ok {
		def.SetDefaults()
	}

	if errDetail, err := validator.ValidateStruct(data); err != nil {
		detail := any(err)
		if errDetail != nil {
			detail = errDetail
		}
		return response.Error(c, fiber.StatusBadRequest, source+" validation error", response.ErrValidation, detail)
	}

	var parsed any
	var err error
	if len(params) > 0 {
		if pp, ok := any(&data).(ParamParser); ok {
			parsed, err = pp.ValidateAndParseWithParams(params[0])
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, "Validation failed", response.ErrValidation, err)
			}
		}
	}

	if parsed == nil {
		if p, ok := any(&data).(Parser); ok {
			parsed, err = p.ValidateAndParse()
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, "Validation failed", response.ErrValidation, err)
			}
		}
	}

	if parsed != nil {
		c.Locals(localKey, parsed)
	} else {
		c.Locals(localKey, data)
	}
	if storeRaw && parsed != nil {
		c.Locals(localKey+"_raw", data)
	}

	return c.Next()
}

func ValidateBody[T any](storeRaw bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "body", c.BodyParser, "body", storeRaw)
	}
}

func ValidateBodyWithParams[T any](storeRaw bool, params any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "body", c.BodyParser, "body", storeRaw, params)
	}
}

func ValidateQuery[T any](storeRaw bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "query", c.QueryParser, "query", storeRaw)
	}
}

func ValidateParams[T any](storeRaw bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return validate[T](c, "params", c.ParamsParser, "params", storeRaw)
	}
}
