package transaction

import (
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/middlewares"
	"ChainServer/internal/common/types"

	"github.com/gofiber/fiber/v2"
)

type TransactionRoutes struct {
	handler          *TransactionHandler
	transactionGroup fiber.Router
}

func NewTransactionRoutes(dbRepo DbTransactionRepository, utxoRepo DbUTXORepository) *TransactionRoutes {
	service := NewTransactionService(dbRepo, utxoRepo)
	handler := NewTransactionHandler(service)

	return &TransactionRoutes{handler: handler}
}

func (r *TransactionRoutes) InitRoutes(router fiber.Router) {
	r.transactionGroup = router.Group("/txs")
}

func (r *TransactionRoutes) RegisterPublic(router fiber.Router) {
	publicGroup := r.transactionGroup.Group("/__pub")

	publicGroup.Get("/",
		middlewares.ValidateQuery[dto.PaginationQuery](false),
		r.handler.GetListTransaction,
	)

	publicGroup.Get("/search",
		middlewares.ValidateQuery[GetTransactionSearchDto](false),
		r.handler.SearchTransactions,
	)

}

func (r *TransactionRoutes) RegisterPrivate(router fiber.Router) {
	privateGroup := r.transactionGroup.Group("/__pri",
		middlewares.JWTAuthMiddleware[types.JWTWalletAuthPayload],
	)

	privateGroup.Get("/pending",
		middlewares.ValidateQuery[dto.PaginationQuery](false),
		r.handler.GetListTransactionPending,
	)

	privateGroup.Post("/new",
		middlewares.DecryptBodyMiddleware(nil),
		middlewares.ValidateBody[NewTransactionDto](true),
		VerifyCreateSignatureMiddleware,
		r.handler.CreateNewTransaction,
	)

	privateGroup.Post("/send",
		middlewares.DecryptBodyMiddleware(nil),
		middlewares.ValidateBody[SendTransactionDto](true),
		VerifySendPayloadMiddleware,
		r.handler.SendTransaction,
	)

}
