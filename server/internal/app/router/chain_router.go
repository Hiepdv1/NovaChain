package router

import (
	"ChainServer/bootstrap"
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/config"

	"github.com/gofiber/fiber/v2"
)

func RegisterChainRoutes(router fiber.Router) {
	envConfig := bootstrap.AppEnv()

	rpcRepo := chain.NewRPCChainRepository(envConfig)

	dbRepo := chain.NewDBChainRepository(config.DB, envConfig)
	tranRepo := transaction.NewDbTransactionRepository(config.DB)

	service := chain.NewChainService(rpcRepo, dbRepo, tranRepo)
	handler := chain.NewChainHandler(service)

	chain.RegisterRoutes(router, handler)

}
