package transaction

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service *TransactionService
}

func NewTransactionHandler(service *TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (s *TransactionHandler) GetListTransaction(c *fiber.Ctx) error {
	queries := c.Locals("query").(dto.PaginationQuery)

	txs, pagination, appErr := s.service.GetListTransaction(queries)

	if appErr != nil {
		return helpers.HandleAppError(c, appErr)
	}

	return response.SuccessList(
		c,
		txs,
		*pagination,
		"Get list transactions sucessfully",
		fiber.StatusOK,
	)
}
