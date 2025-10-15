package helpers

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Panicf("[PANIC RECOVERED] error=%v\nstack:\n%s", r, string(debug.Stack()))
			}
		}()
		fn()
	}()
}

func RecoverAndLog() {
	if r := recover(); r != nil {
		log.Errorf("[PANIC RECOVERED] error=%v\nstack:\n%s", r, string(debug.Stack()))
	}
}
