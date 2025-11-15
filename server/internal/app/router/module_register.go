package router

import (
	"ChainServer/internal/app/module/applog"
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/dashboard"
	"ChainServer/internal/app/module/download"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/app/module/utxo"
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
			utxo.NewDbUTXORepository(),
		),

		wallet.NewWalletRoutes(
			wallet.NewRPCWalletRepository(),
			wallet.NewDBWalletRepository(),
		),

		utxo.NewUTXORoutes(
			utxo.NewRPCUtxoRepository(),
			utxo.NewDbUTXORepository(),
		),

		dashboard.NewDashboardRoutes(
			chain.NewDBChainRepository(),
			transaction.NewDbTransactionRepository(),
		),

		download.NewDownloadRoutes(),
	}
}
