package dashboard

import "github.com/gofiber/fiber/v2"

type DashboardRoutes struct {
	handler        *DashboardHandler
	dashboardGroup fiber.Router
}

func NewDashboardRoutes(
	chainRepo ChainRepository,
	tranRepo TXRepository,
) *DashboardRoutes {

	service := NewDashboardService(
		tranRepo,
		chainRepo,
	)
	handler := NewDashboardHandler(service)

	return &DashboardRoutes{
		handler: handler,
	}
}

func (r *DashboardRoutes) InitRoutes(router fiber.Router) {
	r.dashboardGroup = router.Group("/dashboard")
}

func (r *DashboardRoutes) RegisterPublic(router fiber.Router) {
	r.dashboardGroup.Get("/",
		r.handler.GetNetworkOverview,
	)

	r.dashboardGroup.Get("/recent-activity",
		r.handler.GetRecentActivity,
	)

}
