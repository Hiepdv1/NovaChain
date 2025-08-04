package applog

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, handler *AppLogHandler) {
	appLogGroup := router.Group("/applog")

	appLogGroup.Get("/error", middlewares.ValidateQuery[dto.PaginationQuery](), handler.GetLogError)
	appLogGroup.Get("/:trace_id", handler.GetLogByTraceID)
}
