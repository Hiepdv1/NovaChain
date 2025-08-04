package utils

import (
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func RecoverAppPanic() {
	if r := recover(); r != nil {
		stack := make([]byte, 8<<10)
		length := runtime.Stack(stack, true)
		stackStrace := string(stack[:length])

		entry := log.WithFields(log.Fields{
			"log_scope": "application",
			"recover":   r,
			"stack":     stackStrace,
			"time":      time.Now(),
		})

		entry.Fatal("💥 Application panic — shutting down...")

		os.Exit(1)
	}
}
