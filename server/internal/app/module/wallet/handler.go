package wallet

import (
	"ChainServer/internal/common/constants"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/types"
	"ChainServer/internal/common/utils"
	"encoding/hex"

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

func (h *WalletHandler) CreateWallet(c *fiber.Ctx) error {
	dto := c.Locals("body").(*dto.WalletParsed)

	token, apperr := h.service.CreateWallet(dto)

	if apperr != nil {
		return apperr.Response(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieAccessToken,
		Value:    *token,
		Path:     "/",
		Domain:   env.Cfg.Domain_Client,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "strict",
		// MaxAge:      int(env.Cfg.Jwt_TTL_Minutes * 60),
		SessionOnly: true,
	})

	return response.Success(
		c,
		nil,
		"Create wallet successfully",
		fiber.StatusCreated,
	)
}

func (h *WalletHandler) ImportWallet(c *fiber.Ctx) error {
	dto := c.Locals("body").(*dto.WalletParsed)

	token, apperr := h.service.ImportWallet(dto)

	if apperr != nil {
		return apperr.Response(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     constants.CookieAccessToken,
		Value:    *token,
		Path:     "/",
		Domain:   env.Cfg.Domain_Client,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "strict",
		// MaxAge:      int(env.Cfg.Jwt_TTL_Minutes * 60),
		SessionOnly: true,
	})

	return response.Success(
		c,
		nil,
		"Import wallet successfully",
		fiber.StatusOK,
	)
}

func (h *WalletHandler) Disconnect(c *fiber.Ctx) error {
	payload, ok := c.Locals("wallet").(*utils.JWTPayload[types.JWTWalletAuthPayload])
	token := c.Cookies(constants.CookieAccessToken)

	if !ok || token == "" {
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			"Something went wrong. Please try again.",
			response.ErrInternal,
			nil,
		)
	}

	h.service.Disconnect(token, *payload)

	c.ClearCookie(constants.CookieAccessToken, "/", env.Cfg.Domain_Client)

	return response.Success(
		c,
		nil,
		"wallet disconnected successfully",
		fiber.StatusOK,
	)

}

func (h *WalletHandler) GetMe(c *fiber.Ctx) error {
	payload := c.Locals("wallet").(*utils.JWTPayload[types.JWTWalletAuthPayload])

	pubkey, err := hex.DecodeString(payload.Data.Pubkey)

	if err != nil {
		return response.Error(
			c,
			fiber.StatusInternalServerError,
			"Something went wrong. Please try again.",
			response.ErrInternal,
			nil,
		)
	}

	wallet, apperr := h.service.GetWallet(pubkey)

	if apperr != nil {
		return apperr.Response(c)
	}

	return response.Success(
		c,
		wallet,
		"Get wallet information successfully",
		fiber.StatusOK,
	)
}

func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {

	return response.Success(
		c,
		0,
		"Get balance successfully",
		fiber.StatusOK,
	)
}
