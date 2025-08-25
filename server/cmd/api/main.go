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

	if err := app.Listen(fmt.Sprint(":", env.Cfg.ServerPort)); err != nil {
		log.Error("Listening server error: ", err)
	}
}
