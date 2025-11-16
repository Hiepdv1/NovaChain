package download

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	"os"

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

	queries, appErr := helpers.GetLocalQuery[DowloadFileQuery](c)
	if appErr != nil {
		return appErr.Response(c)
	}

	filePath, appErr := h.service.DowloadNovachain(params.Filename)

	if appErr != nil {
		return appErr.Response(c)
	}

	file, err := os.Stat(filePath)
	if err != nil {
		return apperror.NotFound("File not found", nil).Response(c)
	}

	if queries.Query == "info" {
		return response.Success(
			c,
			DownloadInfo{
				Name: params.Filename,
				Size: file.Size(),
			},
			"Download Info Retrieved Successfully",
			fiber.StatusOK,
		)
	}

	return response.SendFileCustom(
		c,
		filePath,
		params.Filename,
	)
}
