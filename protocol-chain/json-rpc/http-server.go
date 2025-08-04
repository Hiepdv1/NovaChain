package jsonrpc

import (
	"core-blockchain/cmd/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func StartHTTPServer(port, addr string, rpcEnable bool, cli *utils.CommandLine) {
	api := NewAPI(cli)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "🚀 JSON-RPC 2.0 HTTP Server Running")
	})
	http.HandleFunc("/__jsonrpc", api.HandleHTTPJSONRPC)

	addrPort := fmt.Sprintf("%s:%s", addr, port)
	log.Info("🌐 Serving JSON-RPC over HTTP at ", addrPort)
	http.ListenAndServe(addrPort, nil)
}

func (api *API) HandleHTTPJSONRPC(w http.ResponseWriter, r *http.Request) {
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res := api.ProcessRequest(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
