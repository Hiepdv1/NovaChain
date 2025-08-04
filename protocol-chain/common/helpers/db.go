package helpers

import (
	blockchain "core-blockchain/core"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vrecan/death"
)

func SafeCloseDB(bc *blockchain.Blockchain) {
	if bc == nil || bc.Database == nil {
		log.Warn("Blockchain or database is nil, nothing to close.")
		return
	}

	if err := bc.Database.Close(); err != nil {
		log.Errorf("Failed to close DB: %v", err)
	} else {
		log.Infoln("✅ Database closed successfully.")
	}
}

func SafeGuardDB(bc *blockchain.Blockchain) {
	if r := recover(); r != nil {
		log.Errorf("Recovered from panic: %v", r)
		SafeCloseDB(bc)
		os.Exit(1)
	}
}

func CloseDB(bc *blockchain.Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func() {
		SafeCloseDB(bc)
		os.Exit(0)
		log.Infoln("Shutting down... (SIGINT/SIGTERM)")
	})
}
