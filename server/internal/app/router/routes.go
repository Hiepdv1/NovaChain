package router

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/v1/explorer")

	for _, r := range GetAllModuleRouters() {
		r.InitRoutes(api)

		if publicRouter, ok := r.(PublicRouter); ok {
			publicRouter.RegisterPublic(api)

		}

		if privateRouter, ok := r.(PrivateRouter); ok {
			privateRouter.RegisterPrivate(api)
		}
	}
}
