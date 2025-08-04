package router

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/v1/explorer")

	RegisterAppLogRoutes(api)
	RegisterWalletRoutes(api)
	RegisterChainRoutes(api)
	RegisterTransactionRoutes(api)
}
