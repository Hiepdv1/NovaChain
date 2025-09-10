package utxo

import (
	"github.com/gofiber/fiber/v2"
)

type UTXORoutes struct {
	UTXOGroup fiber.Router
	handler   *UTXOHandler
}

func NewUTXORoutes(
	rpcRepo RPCUtxoRepository,
	dbRepo DbUTXORepository,
) *UTXORoutes {

	service := NewUTXOService(rpcRepo, dbRepo)
	handler := NewUTXOHandler(service)

	return &UTXORoutes{
		handler: handler,
	}
}

func (r *UTXORoutes) InitRoutes(router fiber.Router) {
	r.UTXOGroup = router.Group("/utxos")
}

func (r *UTXORoutes) RegisterPrivate(router fiber.Router) {
	// privateGroup := r.UTXOGroup.Group("/__pri",
	// 	middlewares.JWTAuthMiddleware[types.JWTWalletAuthPayload],
	// )

}
