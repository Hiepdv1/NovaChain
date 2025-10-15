package scheduler

import (
	blocksync "ChainServer/internal/scheduler/jobs/block-sync"
	sendtx "ChainServer/internal/scheduler/jobs/send-tx"

	log "github.com/sirupsen/logrus"
)

func StartSchedulers() {
	log.Info("✅ Schedulers started")
	go blocksync.Run()
	go sendtx.Run()
}
