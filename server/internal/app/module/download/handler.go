package download

import (
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type DownloadHandler struct {
	service *DownloadService
}

func NewDownloadHandler(service *DownloadService) *DownloadHandler {
	return &DownloadHandler{
		service: service,
	}
}

func (h *DownloadHandler) DowloadNovachainHandler(c *fiber.Ctx) error {
	params, appErr := helpers.GetLocalParams[DowloadFileParams](c)
	if appErr != nil {
		return appErr.Response(c)
	}
	filePath, appErr := h.service.DowloadNovachain(params.Filename)
	if appErr != nil {
		return appErr.Response(c)
	}

	return response.SendFileCustom(
		c,
		filePath,
		params.Filename,
	)
}
