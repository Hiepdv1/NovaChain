package scheduler

import (
	"ChainServer/internal/scheduler/jobs/blocksync"

	log "github.com/sirupsen/logrus"
)

func StartSchedulers() {
	log.Info("✅ Schedulers started")
	go blocksync.Run()
}
