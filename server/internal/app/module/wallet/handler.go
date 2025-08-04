package wallet

import (
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	service *WalletService
}

func NewWalletHandler(service *WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {

	return response.Success(
		c,
		0,
		"Get balance successfully",
		fiber.StatusOK,
	)
}
