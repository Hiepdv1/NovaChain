package sendtx

import (
	"ChainServer/internal/app/module/transaction"
	"time"

	log "github.com/sirupsen/logrus"
)

type JobSendTx interface {
	Start()
}

type jobSendTx struct {
	interval time.Duration
	dbTrans  transaction.DbTransactionRepository
	rpcTrans transaction.RpcTransactionRepository
}

func NewJobSendTx(
	interval time.Duration,
	dbTrans transaction.DbTransactionRepository,
	rpcTrans transaction.RpcTransactionRepository,
) JobSendTx {
	return &jobSendTx{
		interval: interval,
		dbTrans:  dbTrans,
		rpcTrans: rpcTrans,
	}
}

func (j *jobSendTx) Start() {
	log.Info("🔄 Send tx job started")
	j.StartSendTx()
}
