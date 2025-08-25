package response

import (
	"github.com/gofiber/fiber/v2"
)

func Success(c *fiber.Ctx, data any, message string, statusCode int) error {
	traceID := GetTraceID(c)

	res := ResponseBody{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		TraceID:    traceID,
	}

	return c.Status(statusCode).JSON(res)
}

func Error(c *fiber.Ctx, statusCode int, rawMsg string, errType ErrorType, err any) error {
	traceID := GetTraceID(c)

	res := ResponseBody{
		Success:    false,
		StatusCode: statusCode,
		Message:    GetMessage(errType, rawMsg),
		Error:      err,
		TraceID:    traceID,
	}

	return c.Status(statusCode).JSON(res)

}

func SuccessList(c *fiber.Ctx, data any, meta PaginationMeta, message string, statusCode int) error {
	traceID := GetTraceID(c)

	res := ListResponse{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       meta,
		TraceID:    traceID,
	}

	return c.Status(statusCode).JSON(res)
}
