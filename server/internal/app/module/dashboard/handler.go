package dashboard

import (
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service *DashboardService
}

func NewDashboardHandler(service *DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) GetNetworkOverview(c *fiber.Ctx) error {

	networkOverview, apperr := h.service.GetNetworkOverview()
	if apperr != nil {
		return apperr.Response(c)
	}

	return response.Success(
		c,
		networkOverview,
		"Get Network Overview Successfully",
		fiber.StatusOK,
	)
}

func (h *DashboardHandler) GetRecentActivity(c *fiber.Ctx) error {

	recentActivity, apperr := h.service.GetRecentActivity()
	if apperr != nil {
		return apperr.Response(c)
	}

	return response.Success(
		c,
		recentActivity,
		"Get Recent Activity Successfully",
		fiber.StatusOK,
	)
}
