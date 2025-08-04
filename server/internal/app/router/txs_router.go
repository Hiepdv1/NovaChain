package router

import (
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/config"

	"github.com/gofiber/fiber/v2"
)

func RegisterTransactionRoutes(router fiber.Router) {

	dbRepo := transaction.NewDbTransactionRepository(config.DB)

	service := transaction.NewTransactionService(dbRepo)
	handler := transaction.NewTransactionHandler(service)

	transaction.RegisterRoutes(router, handler)
}
