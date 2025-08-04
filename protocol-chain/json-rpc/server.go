package jsonrpc

import (
	"context"
	"core-blockchain/cmd/utils"
	"core-blockchain/common/env"
	"core-blockchain/common/helpers"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	SERVER_PORT = env.GetEnvAsStr("SERVER_PORT", "9000")
)

func StartServer(cli *utils.CommandLine, rpcEnable bool, rpcPort, rpcAddr, rpcMode string) {
	if rpcPort != "" {
		SERVER_PORT = rpcPort
	}

	defer cli.Blockchain.Database.Close()
	go helpers.CloseDB(cli.Blockchain)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch rpcMode {
	case "tcp":
		log.Info("Starting TCP server...")
		go StartTCPServer(SERVER_PORT, rpcAddr, cli)
	case "http":
		log.Info("Starting HTTP server...")
		go StartHTTPServer(SERVER_PORT, rpcAddr, rpcEnable, cli)
	case "both":
		log.Info("Starting HTTP and TCP server...")
		go StartHTTPServer(SERVER_PORT, rpcAddr, rpcEnable, cli)
		go StartTCPServer(SERVER_PORT, rpcAddr, cli)
	default:
		log.Fatalf("Invalid rpcMode: %s (use 'tcp', 'http', or 'both')", rpcMode)
	}

	log.Infof("🚀 Server running... Press Ctrl+C to shutdown.")
	<-ctx.Done()
	log.Warn("🛑 Server shutting down...")

	helpers.SafeGuardDB(cli.Blockchain)
}
