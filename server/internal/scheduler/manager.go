package scheduler

import (
	blocksync "ChainServer/internal/scheduler/jobs/block-sync"

	log "github.com/sirupsen/logrus"
)

func StartSchedulers() {
	log.Info("✅ Schedulers started")
	go blocksync.Run()
}
