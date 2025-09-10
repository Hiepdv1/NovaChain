package wallet

import (
	"ChainServer/internal/cache/redis"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"
	"ChainServer/internal/common/types"

	"github.com/gofiber/fiber/v2"
)

type WalletRoutes struct {
	handler     *WalletHandler
	walletGroup fiber.Router
}

func NewWalletRoutes(rpcRepo RPCWalletRepository, dbRepo DBWalletRepository) *WalletRoutes {

	cacheRepo := NewWalletCacheRepository(redis.Client)
	service := NewWalletService(rpcRepo, dbRepo, cacheRepo)
	handler := NewWalletHandler(service)

	return &WalletRoutes{handler: handler}
}

func (r *WalletRoutes) InitRoutes(router fiber.Router) {
	r.walletGroup = router.Group("/wallet")
}

func (r *WalletRoutes) RegisterPublic(router fiber.Router) {
	publicGroup := r.walletGroup.Group("/__pub",
		middlewares.ValidateBody[dto.WalletRequest](),
		middlewares.VerifyWalletSignature,
	)

	publicGroup.Post("/new", r.handler.CreateWallet)
	publicGroup.Post("/import", r.handler.ImportWallet)

}

func (r *WalletRoutes) RegisterPrivate(router fiber.Router) {
	privateGroup := r.walletGroup.Group("/__pri",
		middlewares.JWTAuthMiddleware[types.JWTWalletAuthPayload],
	)

	privateGroup.Get("/me", r.handler.GetMe)
	privateGroup.Post("/disconnect", r.handler.Disconnect)
}
