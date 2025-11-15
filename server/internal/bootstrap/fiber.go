package bootstrap

import (
	"ChainServer/internal/app/router"
	"ChainServer/internal/common/config"
	"ChainServer/internal/common/middlewares"
	"ChainServer/internal/common/ratelimiter"
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

func InitRouter() *fiber.App {
	globalLimiter, err := ratelimiter.NewTokenBucketRateLimiter(ratelimiter.Config{
		Rate:  1000,
		Burst: 3000,
	})
	if err != nil {
		log.Panicf("Failed to initialize global rate limiter: %v", err)
	}

	app := fiber.New(fiber.Config{})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://novaexplorer.netlify.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(middlewares.MiddlewareRecover())

	app.Use(cors.New(config.CorsConfig()))

	app.Use(middlewares.RateLimitMiddleware(middlewares.RateLimitConfig{
		GlobalLimiter: globalLimiter,
		RouteLimiters: map[string]ratelimiter.RateLimiter{},
	}))

	app.Use(middlewares.MiddlewareRequestLogger())

	router.RegisterRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return response.Error(
			c,
			fiber.StatusNotFound,
			"Endpoint not found",
			response.ErrNotFound,
			"Not Found",
		)
	})

	return app
}
