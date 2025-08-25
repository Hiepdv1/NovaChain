package router

import (
	"ChainServer/internal/app/module/applog"
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/app/module/wallet"
)

func GetAllModuleRouters() []ModuleRouter {
	return []ModuleRouter{
		chain.NewChainRoutes(
			chain.NewRPCChainRepository(),
			chain.NewDBChainRepository(),
			transaction.NewDbTransactionRepository(),
		),

		applog.NewAppLogRoutes(
			applog.NewFileAppLogRepository(),
		),

		transaction.NewTransactionRoutes(
			transaction.NewDbTransactionRepository(),
		),

		wallet.NewWalletRoutes(
			wallet.NewRPCWalletRepository(),
			wallet.NewDBWalletRepository(),
		),
	}
}
