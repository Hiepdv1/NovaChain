package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/app/module/wallet"
	"time"
)

func Run() {

	blockSync := NewJobBlockSync(
		time.Minute,
		chain.NewDBChainRepository(),
		transaction.NewDbTransactionRepository(),
		chain.NewRPCChainRepository(),
		wallet.NewDBWalletRepository(),
	)

	blockSync.Start()
}
