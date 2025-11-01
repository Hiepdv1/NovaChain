package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/app/module/utxo"
	"ChainServer/internal/app/module/wallet"
	"time"
)

func Run() {

	blockSync := NewJobBlockSync(
		time.Second,
		chain.NewDBChainRepository(),
		transaction.NewDbTransactionRepository(),
		chain.NewRPCChainRepository(),
		wallet.NewDBWalletRepository(),
		utxo.NewDbUTXORepository(),
	)

	blockSync.Start()
}
