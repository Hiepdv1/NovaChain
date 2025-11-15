package response

import (
	"fmt"
	"os"

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

func SendFileCustom(c *fiber.Ctx, filePath string, fileName string) error {
	file, err := os.Stat(filePath)
	if err != nil {
		return Error(c,
			fiber.StatusNotFound,
			"File not found",
			ErrNotFound,
			nil,
		)
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Set("Content-Length", fmt.Sprintf("%d", file.Size()))

	if c.Method() == "head" {
		return nil
	}

	return c.SendFile(filePath)
}
