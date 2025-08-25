package helpers

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

func HandleAppError(c *fiber.Ctx, appErr *apperror.AppError) error {
	return response.Error(
		c,
		int(appErr.Status),
		appErr.Error(),
		appErr.ErrType,
		appErr.Err,
	)
}
