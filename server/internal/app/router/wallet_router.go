package router

import (
	"ChainServer/bootstrap"
	"ChainServer/internal/app/module/wallet"
	"ChainServer/internal/common/config"

	"github.com/gofiber/fiber/v2"
)

func RegisterWalletRoutes(router fiber.Router) {

	envConfig := bootstrap.AppEnv()

	rpcRepo := wallet.NewRPCWalletRepository(envConfig.Fullnode_RPC_URL)
	dbRepo := wallet.NewDBWalletRepository(config.DB)

	service := wallet.NewWalletService(rpcRepo, dbRepo)
	handler := wallet.NewWalletHandler(service)

	wallet.RegisterRoutes(router, handler)
}
