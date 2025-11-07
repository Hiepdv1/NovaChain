package main

import (
	utilCmd "core-blockchain/cmd/utils"
	"core-blockchain/common/env"
	utilLog "core-blockchain/common/utils"
	blockchain "core-blockchain/core"
	jsonrpc "core-blockchain/json-rpc"
	"core-blockchain/p2p"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Entry point
func main() {
	defer os.Exit(0)

	conf := env.New()

	// Global flags
	var (
		address    string
		instanceID string
		rpcPort    string
		rpcAddress string
		rpcMode    string
		rpcEnabled bool
		isSeedPeer bool
	)

	cli := utilCmd.CommandLine{
		Blockchain: &blockchain.Blockchain{
			Database:   nil,
			InstanceId: instanceID,
		},
		P2P: nil,
	}

	//------------------------------------------------------------
	// INIT COMMAND
	//------------------------------------------------------------
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the blockchain and create the genesis block",
		Long: `Initialize a new blockchain instance with a genesis block.

Required:
  --InstanceId   Unique ID for this blockchain instance (auto-generated if omitted).

Example:
  blockchain init --InstanceId 1001
`,
		Run: func(cmd *cobra.Command, args []string) {
			if instanceID == "" {
				instanceID = fmt.Sprintf("%d", time.Now().Unix())
				log.Infof("No instance ID provided. Generated default ID: %s", instanceID)
			}

			cli, err := cli.UpdateInstance(instanceID, false)
			if err != nil {
				log.Fatal(err)
			}

			cli.CreateBlockchain()
			log.Infof("✅ Blockchain initialized successfully (Instance ID: %s)", instanceID)
			cli.Blockchain.Database.Close()
		},
	}

	//------------------------------------------------------------
	// WALLET COMMAND GROUP
	//------------------------------------------------------------
	walletCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Manage blockchain wallets",
		Long: `Create, list, or check wallet balances.

Subcommands:
  - new       Create a new wallet.
  - list      List all wallets.
  - balance   Get the balance of a wallet address.
`,
	}

	// Create new wallet
	newWalletCmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new wallet",
		Long:  `Generate a new wallet and print its address.`,
		Run: func(cmd *cobra.Command, args []string) {
			cli.CreateWallet()
		},
	}

	// List wallets
	listWalletCmd := &cobra.Command{
		Use:   "list",
		Short: "List all wallet addresses",
		Long:  `List all wallet addresses stored locally.`,
		Run: func(cmd *cobra.Command, args []string) {
			cli.ListWallet()
		},
	}

	// Check wallet balance
	walletBalanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Check a wallet's balance",
		Long: `Get the current balance of a given wallet address.

Required:
  --Address     Wallet address to check (must be 34 characters).
  --InstanceId  Blockchain instance ID.

Example:
  blockchain wallet balance --Address <wallet_address> --InstanceId 1001
`,
		Run: func(cmd *cobra.Command, args []string) {
			if address == "" {
				log.Fatal("--Address flag is required")
			}
			if instanceID == "" {
				log.Fatal("--InstanceId flag is required")
			}

			cli, err := cli.UpdateInstance(instanceID, true)
			if err != nil {
				log.Fatal(err)
			}

			cli.GetBalance(address)
		},
	}

	walletCmd.AddCommand(newWalletCmd, listWalletCmd, walletBalanceCmd)

	//------------------------------------------------------------
	// PRINT COMMANDS
	//------------------------------------------------------------
	printCmd := &cobra.Command{
		Use:   "print",
		Short: "Print all blocks in the blockchain",
		Long:  `Display all blocks stored in the blockchain, including their transactions and metadata.`,
		Run: func(cmd *cobra.Command, args []string) {
			if instanceID == "" {
				log.Fatal("--InstanceId flag is required")
			}

			cli, err := cli.UpdateInstance(instanceID, true)
			if err != nil {
				log.Fatal(err)
			}
			cli.PrintBlockchain()
		},
	}

	utxosCmd := &cobra.Command{
		Use:   "utxos",
		Short: "Show current UTXO set",
		Long:  `Display all unspent transaction outputs currently stored in the blockchain database.`,
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := cli.UpdateInstance(instanceID, true)
			if err != nil {
				log.Fatal(err)
			}
			cli.PrintUtxos()
		},
	}

	bestChainCmd := &cobra.Command{
		Use:   "best-chain",
		Short: "Print the best (longest) chain",
		Long:  `Display the blocks in the current best chain.`,
		Run: func(cmd *cobra.Command, args []string) {
			if instanceID == "" {
				log.Fatal("--InstanceId flag is required")
			}

			cli, err := cli.UpdateInstance(instanceID, true)
			if err != nil {
				log.Fatal(err)
			}
			cli.PrintBestChain()
		},
	}

	//------------------------------------------------------------
	// NODE COMMAND
	//------------------------------------------------------------
	var (
		minerAddress string
		miner        bool
		fullNode     bool
		listenPort   string
	)

	nodeCmd := &cobra.Command{
		Use:   "startNode",
		Short: "Start a blockchain node",
		Long: `Start a full node or mining node in the blockchain network.

Required:
  --Port         Port to listen on.
  --InstanceId   Blockchain instance ID.

Optional:
  --Miner        Run as a mining node (true/false).
  --Fullnode     Run as a full node (true/false).
  --Address      Miner address (required if --Miner is true).
  --SeedPeer     Enable seed peer discovery.
  --RPC          Enable JSON-RPC server.
  --RPC-Port     RPC server port (default: 9000).
  --RPC-Addr     RPC server bind address (default: 127.0.0.1).
  --RPC-Mode     Choose 'http', 'tcp', or 'both'.

Example:
  blockchain startNode --Port 3000 --InstanceId 1001 --Miner true --Address <wallet_address>
`,
		Run: func(cmd *cobra.Command, args []string) {
			if instanceID == "" {
				log.Fatal("--InstanceId flag is required")
			}
			if listenPort == "" {
				log.Fatal("--Port flag is required")
			}
			if miner && minerAddress == "" {
				log.Fatal("Miner address is required when --Miner is true")
			}

			logPath := fmt.Sprintf("/logs/console_%s.log", instanceID)
			_ = utilLog.ClearLogFile(logPath)

			cli, err := cli.UpdateInstance(instanceID, false)
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("🔧 Starting node (Instance: %s, Port: %s)", instanceID, listenPort)

			cli.StartNode(listenPort, minerAddress, miner, fullNode, isSeedPeer, func(net *p2p.Network) {
				log.Info("✅ Node started successfully")
				if miner {
					log.Info("Running in Miner mode")
				}
				if rpcEnabled {
					cli.P2P = net
					go jsonrpc.StartServer(cli, rpcEnabled, rpcPort, rpcAddress, rpcMode)
				}
			})
		},
	}

	nodeCmd.Flags().StringVar(&listenPort, "Port", conf.ListenPort, "Node listening port (Required)")
	nodeCmd.Flags().StringVar(&minerAddress, "Address", conf.MinerAddress, "Miner address (Required if --Miner=true)")
	nodeCmd.Flags().BoolVar(&miner, "Miner", conf.Miner, "Set true to enable mining mode")
	nodeCmd.Flags().BoolVar(&fullNode, "Fullnode", conf.FullNode, "Set true to run as a full node")
	nodeCmd.Flags().BoolVar(&isSeedPeer, "SeedPeer", false, "Enable seed peer discovery")

	//------------------------------------------------------------
	// SEND COMMAND
	//------------------------------------------------------------
	var (
		mineNow  bool
		sendFrom string
		sendTo   string
		amount   float64
		fee      float64
	)

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send tokens between wallets",
		Long: `Send tokens from one wallet to another.

Required:
  --SendFrom   Sender wallet address (34 characters)
  --SendTo     Recipient wallet address (34 characters)
  --Amount     Amount to send (> 0)

Optional:
  --Fee        Transaction fee (default: 1/PER_COIN)
  --Mine       Immediately mine the transaction

Example:
  blockchain send --SendFrom <sender> --SendTo <receiver> --Amount 10 --Fee 0.01 --Mine
`,
		Run: func(cmd *cobra.Command, args []string) {
			if sendFrom == "" || sendTo == "" || amount <= 0 {
				log.Fatal("Missing required flags: --SendFrom, --SendTo, and --Amount are mandatory")
			}
			// cli.Send(sendFrom, sendTo, amount, fee, mineNow)
		},
	}

	sendCmd.Flags().StringVar(&sendFrom, "SendFrom", "", "Sender's wallet address")
	sendCmd.Flags().StringVar(&sendTo, "SendTo", "", "Receiver's wallet address")
	sendCmd.Flags().Float64Var(&amount, "Amount", 0, "Amount to send (> 0)")
	sendCmd.Flags().Float64Var(&fee, "Fee", float64(1/blockchain.PER_COIN), "Transaction fee")
	sendCmd.Flags().BoolVar(&mineNow, "Mine", false, "Mine transaction immediately")

	//------------------------------------------------------------
	// ROOT COMMAND
	//------------------------------------------------------------
	rootCmd := &cobra.Command{
		Use:   "blockchain",
		Short: "Blockchain CLI tool",
		Long: `A command-line interface for managing and interacting with the blockchain system.

Global Flags:
  --Address      Wallet address (34 characters)
  --InstanceId   Unique blockchain instance ID
  --RPC-Port     RPC server port (default: 9000)
  --RPC-Addr     RPC server address (default: 127.0.0.1)
  --RPC          Enable JSON-RPC server
  --RPC-Mode     Choose 'http', 'tcp', or 'both'

Example Workflows:
  1. Initialize a blockchain:
     blockchain init --InstanceId 1001

  2. Start a miner node:
     blockchain startNode --Port 3000 --InstanceId 1001 --Miner true --Address <wallet_address>

  3. Create and list wallets:
     blockchain wallet new
     blockchain wallet list

  4. Check balance:
     blockchain wallet balance --Address <wallet_address> --InstanceId 1001

  5. Send tokens:
     blockchain send --SendFrom <addr1> --SendTo <addr2> --Amount 5 --Fee 0.01 --Mine

  6. Inspect the chain:
     blockchain print
`,
	}

	// Global persistent flags
	rootCmd.PersistentFlags().StringVar(&address, "Address", "", "Your wallet address")
	rootCmd.PersistentFlags().StringVar(&instanceID, "InstanceId", "", "Blockchain instance ID")
	rootCmd.PersistentFlags().StringVar(&rpcPort, "RPC-Port", "9000", "RPC server port")
	rootCmd.PersistentFlags().StringVar(&rpcAddress, "RPC-Addr", "127.0.0.1", "RPC server bind address")
	rootCmd.PersistentFlags().BoolVar(&rpcEnabled, "RPC", false, "Enable JSON-RPC server")
	rootCmd.PersistentFlags().StringVar(&rpcMode, "RPC-Mode", "http", "Set RPC mode: 'http', 'tcp', or 'both'")

	rootCmd.AddCommand(
		initCmd,
		walletCmd,
		utxosCmd,
		sendCmd,
		printCmd,
		nodeCmd,
		bestChainCmd,
	)

	rootCmd.Execute()
}
