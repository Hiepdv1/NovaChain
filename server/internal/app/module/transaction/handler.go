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

	txsEncryped, apperr := h.service.CreateNewTransaction(walletPayload, *dto)

	if apperr != nil {
		return apperr.Response(c)
	}

	return response.Success(
		c,
		txsEncryped,
		"Create New Transaction Successfully",
		fiber.StatusCreated,
	)

}

func (h *TransactionHandler) SendTransaction(c *fiber.Ctx) error {
	walletPayload, apperr := helpers.GetLocalWallet(c)
	if apperr != nil {
		return apperr.Response(c)
	}

	dto, apperr := helpers.GetLocalBody[SendTransactionDataParsed](c)
	if apperr != nil {
		return apperr.Response(c)
	}

	apperr = h.service.SendTransaction(walletPayload, dto)

	if apperr != nil {
		return apperr.Response(c)
	}

	return response.Success(
		c,
		nil,
		"Send Transaction Successfully",
		fiber.StatusCreated,
	)
}

func (h *TransactionHandler) GetListTransactionPending(c *fiber.Ctx) error {

	wallet, apperr := helpers.GetLocalWallet(c)
	if apperr != nil {
		return apperr.Response(c)
	}

	pagination, apperr := helpers.GetLocalQuery[dto.PaginationQuery](c)
	if apperr != nil {
		return apperr.Response(c)
	}

	txPendings, paginationMeta, apperr := h.service.TransactionPending(wallet, pagination)
	if apperr != nil {
		return apperr.Response(c)
	}

	return response.SuccessList(
		c,
		txPendings,
		*paginationMeta,
		"Get list pending transactions successfully",
		fiber.StatusOK,
	)
}

func (h *TransactionHandler) SearchTransactions(c *fiber.Ctx) error {
	queries, apperr := helpers.GetLocalQuery[GetTransactionSearchDto](c)
	if apperr != nil {
		return apperr.Response(c)
	}

	txs, pagination, appErr := h.service.SearchTransactions(queries)

	if appErr != nil {
		return appErr.Response(c)
	}

	return response.SuccessList(
		c,
		txs,
		*pagination,
		"Search transactions successfully",
		fiber.StatusOK,
	)
}
