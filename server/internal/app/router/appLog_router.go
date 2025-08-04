package router

import (
	"ChainServer/internal/app/module/applog"
	"ChainServer/internal/common/config"

	"github.com/gofiber/fiber/v2"
)

func RegisterAppLogRoutes(router fiber.Router) {
	root := config.AppRoot

	repo := applog.NewFileAppLogRepository(root)
	service := applog.NewAppLogService(repo)
	handler := applog.NewAppLogHandler(service)

	applog.RegisterRoutes(router, handler)
}
