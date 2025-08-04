package response

import (
	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, data any, message string, statusCode int) error {
	traceID := GetTraceID(c)

	res := ResponseBody{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		TraceID:    traceID,
	}

	return c.Status(statusCode).JSON(res)
}

func Error(c *fiber.Ctx, statusCode int, rawMsg string, code ErrorCode, err any, stack any) error {
	traceID := GetTraceID(c)

	res := ResponseBody{
		StatusCode: statusCode,
		Message:    GetMessage(code, rawMsg),
		Error:      err,
		TraceID:    traceID,
	}

	if stack != nil {
		res.Stack = stack
	}

	return c.Status(statusCode).JSON(res)

}

func SuccessList(c *fiber.Ctx, data any, meta PaginationMeta, message string, statusCode int) error {
	traceID := GetTraceID(c)

	res := ListResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       meta,
		TraceID:    traceID,
	}

	return c.Status(statusCode).JSON(res)
}
