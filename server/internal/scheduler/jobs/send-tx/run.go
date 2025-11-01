package sendtx

import (
	"ChainServer/internal/app/module/transaction"
	"time"
)

func Run() {
	sendTx := NewJobSendTx(
		time.Second,
		transaction.NewDbTransactionRepository(),
		transaction.NewRPCTransactionRepo(),
	)

	sendTx.Start()
}
