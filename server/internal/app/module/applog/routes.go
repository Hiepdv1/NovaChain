package applog

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

type AppLogRoutes struct {
	handler     *AppLogHandler
	appLogGroup fiber.Router
}

func NewAppLogRoutes(dbRepo AppLogRepository) *AppLogRoutes {
	service := NewAppLogService(dbRepo)
	handler := NewAppLogHandler(service)

	return &AppLogRoutes{handler: handler}
}

func (r *AppLogRoutes) InitRoutes(router fiber.Router) {
	r.appLogGroup = router.Group("/applog")
}

func (r *AppLogRoutes) RegisterPrivate(router fiber.Router) {
	r.appLogGroup.Get("/error", middlewares.ValidateQuery[dto.PaginationQuery](), r.handler.GetLogError)
	r.appLogGroup.Get("/:trace_id", r.handler.GetLogByTraceID)
}
