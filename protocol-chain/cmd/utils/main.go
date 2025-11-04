package utils

import (
	"core-blockchain/common/env"
	"core-blockchain/common/err"
	"core-blockchain/common/helpers"
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
	defer helpers.RecoverAndLog()
	if blockchain.Exists(cli.Blockchain.InstanceId) {
		log.Infof("Blockchain already exists for instance ID: %s", cli.Blockchain.InstanceId)
		log.Info("Path: ", blockchain.GetDatabasePath(cli.Blockchain.InstanceId))
		return
	}

	chain, err := blockchain.InitBlockchain(cli.Blockchain.InstanceId)
	if err != nil {
		log.Error(err)
		return
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	utxos := blockchain.UTXOSet{Blockchain: chain}
	utxos.Compute()

	log.Info("✅ Initialized Blockchain Successfully")

}

func (cli *CommandLine) StartNode(listenPort, minerAddress string, miner, fullNode, isSeedPeer bool, callback func(*p2p.Network)) {
	defer helpers.RecoverAndLog()
	if miner {
		log.Infof("Starting Node %s as a MINER", listenPort)
		if len(minerAddress) > 0 {
			log.Info("Mining is ON. Address to receive rewards: ", minerAddress)
		} else {
			log.Fatal("Please provide a valid miner address")
		}
	} else {
		log.Infof("Starting Node on PORT: %s", listenPort)
	}

	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		log.Error(err)
		return
	}
	p2p.StartNode(chain, listenPort, minerAddress, miner, fullNode, isSeedPeer, callback)
}

func (cli *CommandLine) UpdateInstance(InstanceId string, closeDbAlways bool) (*CommandLine, error) {
	defer helpers.RecoverAndLog()
	utils.SetLog(InstanceId)
	cli.Blockchain.InstanceId = InstanceId

	if blockchain.Exists(InstanceId) {
		chain, err := cli.Blockchain.ContinueBlockchain()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		cli.Blockchain = chain
	}

	cli.CloseDbAlways = closeDbAlways

	return cli, nil
}

func (cli *CommandLine) SendTx(txs []*blockchain.Transaction) SendResponse {
	defer helpers.RecoverAndLog()
	listTxs := make([]*blockchain.Transaction, 0)
	listTxsStr := make([]string, 0)

	if len(txs) == 0 {
		return SendResponse{
			Message: "Empty",
			Count:   0,
			ListTxs: []string{},
			Error:   nil,
		}
	}

	for _, tx := range txs {
		txID := hex.EncodeToString(tx.ID)

		if !cli.Blockchain.VerifyTransaction(tx) {
			log.Warnf("Verify failed: %s", txID)
			return SendResponse{
				Error: err.ErrInvalidArgument("Transaction invalid", txID),
			}
		}
		if p2p.MemoryPool.HasTX(txID) {
			continue
		}
		listTxs = append(listTxs, tx)
		listTxsStr = append(listTxsStr, txID)
	}

	if cli.P2P != nil {
		if len(listTxs) > 0 {
			cli.P2P.Transactions <- listTxs

		} else {
			return SendResponse{
				Message: "Empty",
				Count:   0,
				ListTxs: []string{},
				Error:   nil,
			}
		}
	} else {
		return SendResponse{
			Error: err.ErrInternal("Internal error"),
		}
	}

	return SendResponse{
		Message: "Send transaction successfully",
		Count:   int64(len(listTxsStr)),
		ListTxs: listTxsStr,
		Error:   nil,
	}
}

func (cli *CommandLine) GetMiningTxs(verbose bool) GetMiningTxsResponse {
	defer helpers.RecoverAndLog()
	listTxs := make([]any, 0)

	if verbose {
		for _, txInfo := range p2p.MemoryPool.Queued {
			listTxs = append(listTxs, txInfo.Transaction)
		}
	} else {
		for _, txInfo := range p2p.MemoryPool.Queued {
			listTxs = append(listTxs, hex.EncodeToString(txInfo.Transaction.ID))
		}
	}

	return GetMiningTxsResponse{
		Message: "Get mining transactions successfully",
		ListTxs: listTxs,
		Count:   int64(len(listTxs)),
		Error:   nil,
	}
}

func (cli *CommandLine) ComputeUTXOs() {
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		log.Error(err)
		return
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	utxos := blockchain.UTXOSet{
		Blockchain: chain,
	}
	utxos.Compute()
	count, err := utxos.CountTransactions()
	if err != nil {
		log.Errorf("Count Transaction with error: %v", err)
		return
	}
	log.Infof("Rebuild DONE!!!, there are %d transactions in the UTXOs set", count)
}

func (cli *CommandLine) PrintUtxos() {
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		log.Error(err)
		return
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	utxos, err := chain.FindUTXO()

	if err != nil {
		log.Errorf("Find utxos with error: %v", err)
		return
	}

	for _, u := range utxos {
		for _, out := range u.Outputs {
			fmt.Printf("Pub_key_Hash: %x\n", out.PubKeyHash)
			fmt.Printf("Value: %f\n", out.Value)

			fmt.Printf("-------------------------------------------\n")
		}
	}

	log.Infof("Rebuild DONE!!!")
}

func (cli *CommandLine) PrintBestChain() {
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		log.Error(err)
		return
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	block, err := chain.GetLastBlock()
	if err != nil {
		log.Errorf("Find last block with error: %v", err)
		return
	}

	fmt.Printf("-------------------------------------------\n\n")

	fmt.Printf("Current_Hash: %x\n", block.Hash)
	fmt.Printf("Prev_Hash: %x\n", block.PrevHash)
	fmt.Printf("Height: %d\n", block.Height)
	fmt.Printf("NchainWork: %s\n", block.NChainWork.String())
	fmt.Printf("NBit: %d\n", block.NBits)
	fmt.Printf("Is_Valid: %v", chain.IsBlockValid(*block))

	fmt.Println("\n-------------------------------------------")

}

func (cli *CommandLine) GetBalance(address string) BalanceResponse {
	if !wallet.ValidateAddress(address) {
		return BalanceResponse{
			Address:   address,
			Timestamp: time.Now().Unix(),
			Error:     err.ErrInvalidArgument("SendFrom address is invalid"),
		}
	}
	chain, e := cli.Blockchain.ContinueBlockchain()
	if e != nil {
		log.Error(e)
		return BalanceResponse{
			Address:   address,
			Timestamp: time.Now().Unix(),
			Error:     err.ErrInternal("Internal error"),
		}
	}
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	balance := float64(0)
	publicKeyHash := wallet.Base58Decode([]byte(address))

	publicKeyHash = publicKeyHash[1 : int64(len(publicKeyHash))-checkSumlength]
	utxos := blockchain.UTXOSet{
		Blockchain: chain,
	}

	UTXOs, e := utxos.FindUnSpentTransactions(publicKeyHash)
	if e != nil {
		log.Errorf("Find UTXOS With Error: %v", e.Error())
		return BalanceResponse{
			Address:   address,
			Timestamp: time.Now().Unix(),
			Error:     err.ErrInternal("Internal error"),
		}
	}
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

func (cli *CommandLine) CreateWallet() (string, error) {
	cwd := false
	wallets, err := wallet.InitializeWallets(cwd)
	if err != nil {
		return "", err
	}
	address := wallets.AddWallet()

	wallets.SaveFile(cwd)

	log.Infof("NEW WALLET WITH ADDRESS: %s", address)

	return address, nil
}

func (cli *CommandLine) ListWallet() error {
	cwd := false
	wallets, err := wallet.InitializeWallets(cwd)
	if err != nil {
		return err
	}

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

	return nil
}

func (cli *CommandLine) PrintBlockchain() {
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		log.Error(err)
		return
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	iter, err := chain.Iterator()
	if err != nil {
		log.Error(err)
		return
	}

	for {
		block, _ := iter.Next()
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

func (cli *CommandLine) GetBlockChain(startHash []byte, max uint16) ([]*blockchain.Block, error) {
	var blocks []*blockchain.Block

	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		return nil, err
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	if len(startHash) > 0 {
		currentHash := startHash
		for {
			block, err := chain.GetBlock(currentHash)
			if err != nil {
				return nil, err
			}

			blocks = append(blocks, &block)

			if len(block.PrevHash) == 0 || len(blocks) == int(max) {
				break
			}
		}

	} else {
		iter, err := chain.Iterator()
		if err != nil {
			return nil, err
		}
		for {
			block, err := iter.Next()
			if err != nil {
				return blocks, nil
			}

			blocks = append(blocks, block)

			if len(block.PrevHash) == 0 || len(blocks) == int(max) {
				break
			}
		}
	}

	return blocks, nil

}

func (cli *CommandLine) GetBlockByHeight(height int64) (*blockchain.Block, error) {
	var block *blockchain.Block
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		return nil, err
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	iter, err := chain.Iterator()
	if err != nil {
		return nil, err
	}
	for {
		block, err := iter.Next()
		if err != nil {
			return nil, err
		}
		if block.Height == height {
			return block, nil
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return block, nil
}

func (cli *CommandLine) GetBlocksByHeightRange(height, max int64) ([]*blockchain.Block, error) {
	var blocks []*blockchain.Block
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		return nil, err
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	iter, err := chain.Iterator()
	if err != nil {
		return nil, err
	}

	for {
		block, err := iter.Next()
		if err != nil {
			return blocks, nil
		}

		blocks = append(blocks, block)

		if len(blocks) > int(max) {
			blocks = append(make([]*blockchain.Block, 0), blocks[len(blocks)-int(max):]...)
		}

		if block.Height == height || len(block.PrevHash) == 0 {
			break
		}
	}

	return blocks, nil
}

func (cli *CommandLine) GetCommonBlock(locator [][]byte) (*blockchain.Block, error) {
	var commonBlock *blockchain.Block
	chain, err := cli.Blockchain.ContinueBlockchain()
	if err != nil {
		return nil, err
	}

	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	for _, bHash := range locator {
		block, err := chain.GetBlockMainChain(bHash)
		if err != nil {
			continue
		}

		commonBlock = &block
		break
	}

	return commonBlock, nil
}

func (cli *CommandLine) GetBlockByHash(hash []byte) GetBlockResponse {
	chain, e := cli.Blockchain.ContinueBlockchain()
	if e != nil {
		log.Error(e)
		return GetBlockResponse{
			Error: err.ErrInternal("Internal error"),
		}
	}
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
	chain, e := cli.Blockchain.ContinueBlockchain()
	if e != nil {
		log.Error(e)
		return &GetAllUTXOsResponse{
			Error: err.ErrInternal("Internal error"),
		}
	}
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	txOutputs, e := chain.FindUTXO()
	if e != nil {
		return &GetAllUTXOsResponse{
			Message: "Successfully",
			Error: &err.RPCError{
				Code:    blockchain.HALVING_INTERVAL,
				Message: "Please try again.",
			},
		}
	}

	return &GetAllUTXOsResponse{
		Message: "Successfully",
		Data:    txOutputs,
		Count:   int64(len(txOutputs)),
		Error:   nil,
	}

}

func (cli *CommandLine) GetLastHeight() (int64, error) {
	chain, e := cli.Blockchain.ContinueBlockchain()
	if e != nil {
		log.Errorf("Failed to continue blockchain: %v", e)
		return 0, e
	}
	if cli.CloseDbAlways {
		defer chain.Database.Close()
	}

	lastHeight, err := chain.GetBestHeight()
	if err != nil {
		return 0, err
	}

	return lastHeight, nil
}
