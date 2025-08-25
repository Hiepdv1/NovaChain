package config

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		MaxAge:           int((24 * time.Hour).Seconds()),
	}
}
