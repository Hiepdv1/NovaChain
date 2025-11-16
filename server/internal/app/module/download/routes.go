package download

import (
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

type DownloadRoutes struct {
	handler       *DownloadHandler
	downloadGroup fiber.Router
}

func NewDownloadRoutes() *DownloadRoutes {
	service := NewDownloadService()
	handler := NewDownloadHandler(service)

	return &DownloadRoutes{
		handler: handler,
	}
}

func (r *DownloadRoutes) InitRoutes(router fiber.Router) {
	r.downloadGroup = router.Group("/download")
}

func (r *DownloadRoutes) RegisterPublic(router fiber.Router) {
	r.downloadGroup.Get(
		"/:filename",
		middlewares.ValidateParams[DowloadFileParams](false),
		middlewares.ValidateQuery[DowloadFileQuery](false),
		r.handler.DowloadNovachainHandler,
	)
}
