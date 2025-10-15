package sendtx

import (
	"ChainServer/internal/app/module/transaction"
	"ChainServer/internal/common/constants"
	dbPendingTx "ChainServer/internal/db/pendingTx"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (j *jobSendTx) sendTransactionPending() error {

	ctx := context.Background()

	args := dbPendingTx.SelectPendingTransactionsParams{
		Limit:  30,
		Offset: 0,
		Status: string(constants.TxStatusPending),
	}

	txPendings, err := j.dbTrans.SelectPendingTransactions(ctx, args, nil)
	if err != nil {
		log.Errorf("Find Pending Tx Error: %v", err)
		return err
	}

	var txs []transaction.Transaction
	for _, txPending := range txPendings {
		var tx transaction.Transaction
		err := json.Unmarshal(txPending.RawTx, &tx)
		if err != nil {
			log.Error(err)
			continue
		}
		txs = append(txs, tx)
	}

	res, err := j.rpcTrans.SendTx(txs)
	if err != nil {
		return err
	}

	if res.Error != nil {
		return fmt.Errorf("%s", res.Error.Message)
	}

	log.Info(res.Message)
	log.Infof("List Transaction Sent: %s", strings.Join(res.ListTxs, "\n"))

	return nil
}

func (j *jobSendTx) StartSendTx() {
	log.Info("🔄 Start sending pending transactions")

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for range ticker.C {
		err := j.sendTransactionPending()
		if err != nil {
			log.Errorf("Send Transaction Pending With Error: %v", err.Error())
			continue
		}

	}
}
