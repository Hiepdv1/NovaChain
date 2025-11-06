package chain

import (
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"

	"github.com/gofiber/fiber/v2"
)

type ChainRoutes struct {
	handler    *ChainHandler
	chainGroup fiber.Router
}

func NewChainRoutes(
	rpcRepo RPCChainRepository,
	dbRepo DBChainRepository,
	tranRepo transaction.DbTransactionRepository,
) *ChainRoutes {
	service := NewChainService(rpcRepo, dbRepo, tranRepo)
	handler := NewChainHandler(service)
	return &ChainRoutes{handler: handler}
}

func (r *ChainRoutes) InitRoutes(router fiber.Router) {
	r.chainGroup = router.Group("/chain")
}

func (r *ChainRoutes) RegisterPublic(router fiber.Router) {
	r.chainGroup.Get("/blocks",
		middlewares.ValidateQuery[dto.PaginationQuery](false),
		r.handler.GetBlocks,
	)

	r.chainGroup.Get("/blocks/:blockHash",
		middlewares.ValidateParams[GetBlockDetailParamDto](false),
		middlewares.ValidateQuery[dto.PaginationQuery](false),
		r.handler.GetBlockDetail,
	)

	r.chainGroup.Get("/search",
		middlewares.ValidateQuery[GetSearchResultDto](false),
		r.handler.GetSearchResult,
	)

	r.chainGroup.Get("/network",
		r.handler.GetNetwork,
	)

	r.chainGroup.Get("/miners",
		middlewares.ValidateQuery[dto.PaginationQuery](false),
		r.handler.GetMiners,
	)
}
