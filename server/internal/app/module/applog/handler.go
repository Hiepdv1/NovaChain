package applog

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type AppLogHandler struct {
	service *AppLogService
}

func NewAppLogHandler(service *AppLogService) *AppLogHandler {
	return &AppLogHandler{service: service}
}

func (h *AppLogHandler) GetLogError(c *fiber.Ctx) error {
	queries := c.Locals("query").(dto.PaginationQuery)

	logs, pagination, apperror := h.service.GetListAppLogError(queries)
	if apperror != nil {
		helpers.HandleAppError(c, apperror)
	}

	return response.SuccessList(
		c,
		logs,
		*pagination,
		"Successfully",
		fiber.StatusOK,
	)
}

func (h *AppLogHandler) GetLogByTraceID(c *fiber.Ctx) error {
	traceID := c.Params("trace_id")

	log, apperror := h.service.GetLogByTraceID(traceID)

	if apperror != nil {
		helpers.HandleAppError(c, apperror)
	}

	if log == nil {
		return response.Error(
			c,
			fiber.StatusNotFound,
			"TraceID not found",
			response.ErrNotFound,
			"Not Found",
			nil,
		)
	}

	return response.Success(
		c,
		log,
		"Get AppLog Successfully",
		fiber.StatusOK,
	)

}
