package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/app/module/wallet"
	"ChainServer/internal/common/config"
	"ChainServer/internal/common/env"
	"time"
)

func Run() {

	envConfig := env.New()

	blockSync := NewJobBlockSync(
		time.Minute,
		chain.NewDBChainRepository(config.DB, envConfig),
		transaction.NewDbTransactionRepository(config.DB),
		chain.NewRPCChainRepository(envConfig),
		wallet.NewDBWalletRepository(config.DB),
	)

	blockSync.Start()
}
