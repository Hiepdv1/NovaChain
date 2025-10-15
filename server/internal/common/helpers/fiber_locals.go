package helpers

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/constants"
	"ChainServer/internal/common/types"
	"ChainServer/internal/common/utils"

	"github.com/gofiber/fiber/v2"
)

func SetLocalWallet(c *fiber.Ctx, payload any) {
	c.Locals(constants.LocalsWallet, payload)
}

func GetLocalWallet(c *fiber.Ctx) (*utils.JWTPayload[types.JWTWalletAuthPayload], *apperror.AppError) {
	v := c.Locals(constants.LocalsWallet)
	if v == nil {
		return nil, apperror.Internal("Wallet payload not found in context", nil)
	}

	payload, ok := v.(*utils.JWTPayload[types.JWTWalletAuthPayload])
	if !ok {
		return nil, apperror.Internal("Invalid wallet payload type", nil)
	}

	return payload, nil
}

func GetLocalBody[T any](c *fiber.Ctx) (*T, *apperror.AppError) {
	dto, ok := c.Locals("body").(T)
	if !ok {
		return nil, apperror.Internal("body not found in context", nil)
	}

	return &dto, nil
}

func GetLocalQuery[T any](c *fiber.Ctx) (*T, *apperror.AppError) {
	dto, ok := c.Locals("query").(T)
	if !ok {
		return nil, apperror.Internal("query not found in context", nil)
	}
	return &dto, nil
}

func GetLocalRaw[T any](c *fiber.Ctx, localKey string) (*T, *apperror.AppError) {
	dto, ok := c.Locals(localKey + "_raw").(T)
	if !ok {
		return nil, apperror.Internal("body not found in context", nil)
	}

	return &dto, nil
}
