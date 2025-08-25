package transaction

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

type TransactionRoutes struct {
	handler          *TransactionHandler
	transactionGroup fiber.Router
}

func NewTransactionRoutes(dbRepo DbTransactionRepository) *TransactionRoutes {
	service := NewTransactionService(dbRepo)
	handler := NewTransactionHandler(service)

	return &TransactionRoutes{handler: handler}
}

func (r *TransactionRoutes) InitRoutes(router fiber.Router) {
	r.transactionGroup = router.Group("/txs")
}

func (r *TransactionRoutes) RegisterPublic(router fiber.Router) {
	r.transactionGroup.Get("/", middlewares.ValidateQuery[dto.PaginationQuery](), r.handler.GetListTransaction)
}
