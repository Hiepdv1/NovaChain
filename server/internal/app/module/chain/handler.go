package chain

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type ChainHandler struct {
	service *ChainService
}

func NewChainHandler(service *ChainService) *ChainHandler {
	return &ChainHandler{
		service: service,
	}
}

func (h *ChainHandler) GetBlocks(c *fiber.Ctx) error {

	queries := c.Locals("query").(dto.PaginationQuery)

	blocks, pagination, apperror := h.service.GetBlocks(queries)

	if apperror != nil {
		return helpers.HandleAppError(c, apperror)
	}

	return response.SuccessList(
		c,
		blocks,
		*pagination,
		"Get blocks sucessfully",
		fiber.StatusOK,
	)
}
