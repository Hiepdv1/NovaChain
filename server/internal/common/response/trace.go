package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetTraceID(c *fiber.Ctx) string {
	traceID := c.Get("X-Request-ID")

	if traceID == "" {
		traceID = c.Get("X-Request-Timestamp")
		if traceID != "" {
			return traceID
		}
	}

	if traceID == "" {
		val := c.Locals("requestid")
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}

	return uuid.New().String()
}
