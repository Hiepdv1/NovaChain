package router

import "github.com/gofiber/fiber/v2"

type PublicRouter interface {
	RegisterPublic(router fiber.Router)
}

type PrivateRouter interface {
	RegisterPrivate(router fiber.Router)
}

type ModuleRouter interface {
	InitRoutes(router fiber.Router)
}
