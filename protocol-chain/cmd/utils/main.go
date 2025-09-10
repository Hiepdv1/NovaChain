package utils

import (
	"core-blockchain/common/env"
	"core-blockchain/common/err"
	"core-blockchain/common/utils"
	blockchain "core-blockchain/core"
	"core-blockchain/p2p"
	"core-blockchain/wallet"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

var conf = env.New()

var (
	checkSumlength = conf.WalletAddressCheckSum
)

func (cli *CommandLine) CreateBlockchain() {

	if blockchain.Exists(cli.Blockchain.InstanceId) {
		log.Infof("Blockchain already exists for instance ID: %s", cli.Blockchain.InstanceId)
		log.Info("Path: ", blockchain.GetDatabasePath(cli.Blockchain.InstanceId))
		return
	}

	chain := blockchain.InitBlockchain(cli.Blockchain.InstanceId)

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	utxos := blockchain.UTXOSet{Blockchain: chain}
	utxos.Compute()

	log.Info("✅ Initialized Blockchain Successfully")

}

func (cli *CommandLine) StartNode(listenPort, minerAddress string, miner, fullNode bool, callback func(*p2p.Network)) {
	if miner {
		log.Infof("Starting Node %s as a MINER\n", listenPort)
		if len(minerAddress) > 0 {
			log.Info("Mining is ON. Address to receive rewards: ", minerAddress)
		} else {
			log.Fatal("Please provide a valid miner address")
		}
	} else {
		log.Infof("Starting Node on PORT: %s\n", listenPort)
	}

	chain := cli.Blockchain.ContinueBlockchain()
	p2p.StartNode(chain, listenPort, minerAddress, miner, fullNode, callback)
}

func (cli *CommandLine) UpdateInstance(InstanceId string, closeDbAlways bool) *CommandLine {
	utils.SetLog(InstanceId)
	cli.Blockchain.InstanceId = InstanceId

	if blockchain.Exists(InstanceId) {
		cli.Blockchain = cli.Blockchain.ContinueBlockchain()
	}

	cli.CloseDbAlways = closeDbAlways

	return cli
}

func (cli *CommandLine) SendTx(txs []*blockchain.Transaction) SendResponse {

	listTxs := make([]string, 0)

	for _, tx := range txs {
		if !cli.Blockchain.VerifyTransaction(tx) {
			return SendResponse{
				Error: err.ErrInvalidArgument("Transaction invalid", string(tx.ID)),
			}
		}
		listTxs = append(listTxs, string(tx.ID))
	}

	if cli.P2P != nil {
		for _, tx := range txs {
			cli.P2P.Transactions <- tx
		}
	}

	return SendResponse{
		Message: "Send transaction successfully",
		Count:   int64(len(txs)),
		ListTxs: listTxs,
		Error:   nil,
	}
}

// func (cli *CommandLine) Send(from, to string, amount, fee float64, mineNow bool) SendResponse {
// 	if !wallet.ValidateAddress(from) {
// 		log.Error("SendFrom address is Invalid ")
// 		return SendResponse{
// 			Error: err.ErrInvalidArgument("SendTo address is invalid"),
// 		}
// 	}

// 	if !wallet.ValidateAddress(to) {
// 		log.Error("SendFrom address is Invalid ")
// 		return SendResponse{
// 			Error: err.ErrInvalidArgument("SendFrom address is invalid"),
// 		}
// 	}

// 	chain := cli.Blockchain.ContinueBlockchain()
// 	if cli.CloseDbAlways {
// 		defer chain.Database.Close()
// 	}

// 	utxos := blockchain.UTXOSet{
// 		Blockchain: chain,
// 	}

// tx, err := blockchain.NewTransaction(&wallet, to, amount, fee, &utxos)
// if err != nil {
// 	log.Error(err)
// 	return SendResponse{
// 		Error: &err{
// 			Code:    5028,
// 			Message: "Failed to execute transaction",
// 		},
// 	}
// }

// 	if mineNow {
// 		txs := []*blockchain.Transaction{tx}
// 		log.Info("Transaction executed")

// 		block, err := chain.MineBlock(txs, from, cli.P2P.HandleReoganizeTx)

// 		if err != nil {
// 			log.Error("Failed to mine block: ", err)
// 			return SendResponse{
// 				Error: &Error{
// 					Code:    5028,
// 					Message: "Failed to mine block",
// 				},
// 			}
// 		}

// 		if block == nil {
// 			log.Error("Failed to mine block")
// 			return SendResponse{
// 				SendFrom:  from,
// 				SendTo:    to,
// 				Amount:    amount,
// 				Timestamp: time.Now().Unix(),
// 				Error: &Error{
// 					Code:    5028,
// 					Message: "Failed to mine block",
// 				},
// 			}
// 		}

// 		utxos.Update(block)

// 		if cli.P2P != nil {
// 			cli.P2P.Blocks <- block
// 		}
// 	} else {
// 		if cli.P2P != nil {
// 			cli.P2P.Transactions <- tx
// 		}
// 	}

// 	return SendResponse{
// 		SendFrom:  from,
// 		SendTo:    to,
// 		Amount:    amount,
// 		Timestamp: time.Now().Unix(),
// 		Error:     nil,
// 	}
// }

func (cli *CommandLine) ComputeUTXOs() {
	chain := cli.Blockchain.ContinueBlockchain()

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	utxos := blockchain.UTXOSet{
		Blockchain: chain,
	}
	utxos.Compute()
	count := utxos.CountTransactions()
	log.Infof("Rebuild DONE!!!, there are %d transactions in the UTXOs set", count)
}

func (cli *CommandLine) GetBalance(address string) BalanceResponse {
	if !wallet.ValidateAddress(address) {
		return BalanceResponse{
			Address:   address,
			Timestamp: time.Now().Unix(),
			Error:     err.ErrInvalidArgument("SendFrom address is invalid"),
		}
	}
	chain := cli.Blockchain.ContinueBlockchain()
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	balance := float64(0)
	publicKeyHash := wallet.Base58Decode([]byte(address))

	publicKeyHash = publicKeyHash[1 : int64(len(publicKeyHash))-checkSumlength]
	utxos := blockchain.UTXOSet{
		Blockchain: chain,
	}

	UTXOs := utxos.FindUnSpentTransactions(publicKeyHash)
	for _, out := range UTXOs {
		balance += out.Value
	}

	return BalanceResponse{
		Balance:   balance,
		Address:   address,
		Timestamp: time.Now().Unix(),
		Error:     nil,
	}
}

func (cli *CommandLine) CreateWallet() string {
	cwd := false
	wallets, err := wallet.InitializeWallets(cwd)
	utils.ErrorHandle(err)

	address := wallets.AddWallet()

	wallets.SaveFile(cwd)

	log.Infof("NEW WALLET WITH ADDRESS: %s", address)

	return address
}

func (cli *CommandLine) ListWallet() {
	cwd := false
	wallets, err := wallet.InitializeWallets(cwd)
	utils.ErrorHandle(err)

	addresses := wallets.GetAllAddress()

	for _, address := range addresses {
		fmt.Printf("------------------------------ %s ------------------------------\n", address)
		w, err := wallets.GetWallet(address)
		if err != nil {
			log.Panic("Get wallet with error: ", err)
		}
		privBytes := w.PrivateKey.D.Bytes()
		privHex := hex.EncodeToString(privBytes)
		privWebApp, _ := wallet.EncryptPrivateKeyForExport(privHex)
		fmt.Printf("PrivateKey For WebApp: %s\n", privWebApp)
		fmt.Println()
		fmt.Printf("PublicKey: %s", hex.EncodeToString(w.PublicKey))
		fmt.Println()
	}
}

func (cli *CommandLine) PrintBlockchain() {
	chain := cli.Blockchain.ContinueBlockchain()

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Print("------------------------------------------------\n\n")
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		pow := blockchain.NewProof(block)
		validate := pow.Validate()
		fmt.Printf("Valid: %s\n", strconv.FormatBool(validate))

		fmt.Println("Transactions: ")
		for _, tx := range block.Transactions {
			fmt.Printf(" - %v\n", tx)
		}

		fmt.Print("\n------------------------------------------------\n")

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) GetBlockChain(startHash []byte, max uint16) []*blockchain.Block {
	var blocks []*blockchain.Block

	chain := cli.Blockchain.ContinueBlockchain()

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	if len(startHash) > 0 {
		currentHash := startHash
		for {
			block, err := chain.GetBlock(currentHash)
			utils.ErrorHandle(err)

			blocks = append(blocks, &block)

			if len(block.PrevHash) == 0 || len(blocks) == int(max) {
				break
			}
		}

	} else {
		iter := chain.Iterator()
		for {
			block := iter.Next()

			blocks = append(blocks, block)

			if len(block.PrevHash) == 0 || len(blocks) == int(max) {
				break
			}
		}
	}

	return blocks

}

func (cli *CommandLine) GetBlockByHeight(height int64) blockchain.Block {
	var block blockchain.Block
	chain := cli.Blockchain.ContinueBlockchain()
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	iter := chain.Iterator()
	for {
		block := *iter.Next()
		if block.Height == height {
			return block
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return block
}

func (cli *CommandLine) GetBlocksByHeightRange(height, max int64) []*blockchain.Block {
	var blocks []*blockchain.Block
	chain := cli.Blockchain.ContinueBlockchain()
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	iter := chain.Iterator()

	for {
		block := iter.Next()

		blocks = append(blocks, block)

		if len(blocks) > int(max) {
			blocks = append(make([]*blockchain.Block, 0), blocks[len(blocks)-int(max):]...)
		}

		if block.Height == height || len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks
}

func (cli *CommandLine) GetBlockByHash(hash []byte) GetBlockResponse {
	chain := cli.Blockchain.ContinueBlockchain()
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	block, e := chain.GetBlock(hash)
	if e != nil {
		if errors.Is(e, badger.ErrKeyNotFound) {
			return GetBlockResponse{
				Error: err.ErrNotFound("Block not found"),
			}
		}

		return GetBlockResponse{Error: err.ErrInternal("Get block failed", e.Error())}

	}

	return GetBlockResponse{
		Block: &block,
		Error: nil,
	}
}

func (cli *CommandLine) GetAllUTXOs() *GetAllUTXOsResponse {
	chain := cli.Blockchain.ContinueBlockchain()
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	txOutputs := chain.FindUTXO()

	return &GetAllUTXOsResponse{
		Message: "Successfully",
		Data:    txOutputs,
		Count:   int64(len(txOutputs)),
		Error:   nil,
	}

}
