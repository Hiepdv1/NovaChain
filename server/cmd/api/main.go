package main

import (
	"ChainServer/bootstrap"
	"ChainServer/internal/app/router"
	"ChainServer/internal/common/config"
	"ChainServer/internal/common/middlewares"
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/utils"
	"ChainServer/internal/scheduler"
	"ChainServer/scripts"
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer utils.RecoverAppPanic()

	var envConfig = bootstrap.AppEnv()
	config.InitLogger(envConfig.AppEnv)
	config.InitAppRoot()
	db := config.InitPostgres()
	defer db.Close()
	scripts.AutoMigrate(db)
	scheduler.StartSchedulers()

	app := fiber.New(fiber.Config{})

	app.Use(middlewares.MiddlewareRecover())
	app.Use(middlewares.MiddlewareRequestLogger())

	router.RegisterRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return response.Error(
			c,
			fiber.StatusNotFound,
			"Endpoint not found",
			response.ErrNotFound,
			"Not Found",
			nil,
		)
	})

	if err := app.Listen(fmt.Sprint(":", envConfig.ServerPort)); err != nil {
		log.Error("Listening server error: ", err)
	}
}
