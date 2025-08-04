package jsonrpc

import (
	"core-blockchain/cmd/utils"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	log "github.com/sirupsen/logrus"
)

func StartTCPServer(port, addr string, cli *utils.CommandLine) {
	api := NewAPI(cli)
	rpc.Register(api)

	addrPort := fmt.Sprintf("%s:%s", addr, port)
	ln, err := net.Listen("tcp", addrPort)
	CheckError("Listen TCP error", err)

	log.Info("📡 Serving JSON-RPC over TCP at ", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
