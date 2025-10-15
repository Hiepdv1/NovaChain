package sendtx

import (
	"ChainServer/internal/app/module/transaction"
	"time"
)

func Run() {
	sendTx := NewJobSendTx(
		10*time.Second,
		transaction.NewDbTransactionRepository(),
		transaction.NewRPCTransactionRepo(),
	)

	sendTx.Start()
}
