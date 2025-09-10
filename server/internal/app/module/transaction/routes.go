package transaction

import (
	"ChainServer/internal/app/module/utxo"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"
	"ChainServer/internal/common/types"

	"github.com/gofiber/fiber/v2"
)

type TransactionRoutes struct {
	handler          *TransactionHandler
	transactionGroup fiber.Router
}

func NewTransactionRoutes(dbRepo DbTransactionRepository, utxoRepo utxo.DbUTXORepository) *TransactionRoutes {
	service := NewTransactionService(dbRepo, utxoRepo)
	handler := NewTransactionHandler(service)

	return &TransactionRoutes{handler: handler}
}

func (r *TransactionRoutes) InitRoutes(router fiber.Router) {
	r.transactionGroup = router.Group("/txs")
}

func (r *TransactionRoutes) RegisterPublic(router fiber.Router) {
	publicGroup := r.transactionGroup.Group("/__pub")

	publicGroup.Get("/", middlewares.ValidateQuery[dto.PaginationQuery](), r.handler.GetListTransaction)
}

func (r *TransactionRoutes) RegisterPrivate(router fiber.Router) {
	privateGroup := r.transactionGroup.Group("/__pri",
		middlewares.JWTAuthMiddleware[types.JWTWalletAuthPayload],
	)

	privateGroup.Post("/new",
		VerifyCreateTransactionSig,
		middlewares.ValidateBody[NewTransactionDto](),
		r.handler.CreateNewTransaction,
	)
}
