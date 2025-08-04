package utils

import log "github.com/sirupsen/logrus"

func ErrorHandle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
