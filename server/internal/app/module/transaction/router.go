package transaction

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, handler *TransactionHandler) {
	transactionGroup := router.Group("/txs")

	transactionGroup.Get("/", middlewares.ValidateQuery[dto.PaginationQuery](), handler.GetListTransaction)
}
