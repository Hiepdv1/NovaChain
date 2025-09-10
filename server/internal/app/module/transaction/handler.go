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

func (h *TransactionHandler) GetListTransaction(c *fiber.Ctx) error {
	queries := c.Locals("query").(dto.PaginationQuery)

	txs, pagination, appErr := h.service.GetListTransaction(queries)

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

func (h *TransactionHandler) CreateNewTransaction(c *fiber.Ctx) error {
	dto, apperr := helpers.GetLocalBody[*NewTransactionParsed](c)

	if apperr != nil {
		return apperr.Response(c)
	}

	walletPayload, apperr := helpers.GetLocalWallet(c)

	if apperr != nil {
		return apperr.Response(c)
	}

	tx, apperr := h.service.CreateNewTransaction(walletPayload, *dto)

	if apperr != nil {
		return apperr.Response(c)
	}

	return response.Success(
		c,
		tx,
		"Create New Transaction Successfully",
		fiber.StatusCreated,
	)

}
