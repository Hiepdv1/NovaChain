package chain

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, handler *ChainHandler) {
	chainGroup := router.Group("/chain")

	chainGroup.Get("/blocks", middlewares.ValidateQuery[dto.PaginationQuery](), handler.GetBlocks)
}
