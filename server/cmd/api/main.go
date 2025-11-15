package main

import (
	"ChainServer/internal/bootstrap"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/utils"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func main() {
	defer utils.RecoverAppPanic()

	bootstrap.Init()

	db := bootstrap.StartConnDB()
	defer db.Psql.Close()

	bootstrap.StartSchedulers()

	app := bootstrap.InitRouter()

	certFile := "ssl/cert.pem"
	keyFile := "ssl/key.pem"

	addr := fmt.Sprintf(":%s", env.Cfg.ServerPort)
	if env.Cfg.AppEnv == "development" {
		log.Infof("Starting server on https://localhost%s", addr)
	}

	if err := app.ListenTLS(addr, certFile, keyFile); err != nil {
		log.Error("Listening server error: ", err)
	}
}
