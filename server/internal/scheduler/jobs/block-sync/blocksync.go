package blocksync

import (
	"ChainServer/internal/app/module/chain"
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/app/module/utxo"
	"ChainServer/internal/app/module/wallet"
	"time"

	log "github.com/sirupsen/logrus"
)

type JobBlockSync interface {
	Start()
}

type jobBlockSync struct {
	interval   time.Duration
	dbChain    chain.DBChainRepository
	dbTrans    transaction.DbTransactionRepository
	dbChainRpc chain.RPCChainRepository
	dbWallet   wallet.DBWalletRepository
	dbUtxo     utxo.DbUTXORepository
}

func NewJobBlockSync(
	interval time.Duration,
	dbChain chain.DBChainRepository,
	dbTrans transaction.DbTransactionRepository,
	dbChainRpc chain.RPCChainRepository,
	dbWallet wallet.DBWalletRepository,
	dbUtxo utxo.DbUTXORepository,
) JobBlockSync {
	return &jobBlockSync{
		interval:   interval,
		dbChain:    dbChain,
		dbTrans:    dbTrans,
		dbChainRpc: dbChainRpc,
		dbWallet:   dbWallet,
		dbUtxo:     dbUtxo,
	}
}

func (j *jobBlockSync) Start() {
	log.Info("🔄 Block sync job started with interval:", j.interval)
	j.StartBlockSync(j.interval)
}
