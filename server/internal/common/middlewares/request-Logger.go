package middlewares

import (
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/response"
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func MiddlewareRequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		traceID := uuid.New().String()
		c.Locals("trace_id", traceID)

		err := c.Next()

		stop := time.Since(start)
		entry := log.WithFields(log.Fields{
			"log_scope": "router",
			"time":      time.Now(),
			"method":    c.Method(),
			"path":      c.OriginalURL(),
			"status":    c.Response().StatusCode(),
			"duration":  stop.String(),
			"trace_id":  traceID,
			"ip":        c.IP(),
		})

		status := c.Response().StatusCode()

		switch {
		case err != nil && status >= 500:
			entry.WithField("error", err.Error()).Error("❌ Internal server error")
			return err

		case err != nil:
			entry.WithField("error", err.Error()).Warn("⚠️ Handled error")
			return err

		case status >= 500:
			entry.Error("🔥 Server error.")

		case status >= 400:
			entry.Warn("⚠️ Client error.")

		default:
			entry.Info("✅ Request completed.")
		}

		return nil
	}
}

func MiddlewareRecover() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		traceID := uuid.New().String()
		c.Locals("trace_id", traceID)

		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 8<<10)
				length := runtime.Stack(stack, true)
				stackTrace := string(stack[:length])
				entry := log.WithFields(log.Fields{
					"log_scope": "router",
					"trace_id":  traceID,
					"recover":   r,
					"path":      c.OriginalURL(),
					"ip":        c.IP(),
					"stack":     stackTrace,
					"time":      time.Now(),
				})

				entry.Error("💥 Panic caught")

				if env.Cfg.AppEnv == "production" {

					err = response.Error(
						c,
						fiber.StatusInternalServerError,
						"Internal Server Error",
						response.ErrInternal,
						r,
					)
				} else {
					err = response.Error(
						c,
						fiber.StatusInternalServerError,
						fmt.Sprintf("%v", r),
						response.ErrInternal,
						r,
					)
				}
			}

		}()

		return c.Next()
	}
}
