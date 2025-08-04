package wallet

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *WalletHandler) {
	walletGroup := router.Group("/wallet")

	walletGroup.Get("/:address", handler.GetBalance)
}
