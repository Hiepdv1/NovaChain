package main

import (
	cui "core-blockchain/cmd/demo/CUI"
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

func main() {
	defer os.Exit(0)
	conf := env.New()

	var (
		address    string
		instanceID string
		rpcPort    string
		rpcAddress string
		rpcMode    string
		rpcEnabled bool
		isSeedPeer bool
		chainData  string
		LogFile    string
	)

	cli := utilCmd.CommandLine{
		Blockchain: &blockchain.Blockchain{
			Database:   nil,
			InstanceId: instanceID,
		},
		P2P: nil,
	}

	// -----------------------
	// INIT
	// -----------------------
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize NovaChain blockchain and create the genesis block",
		Long: `Initialize a new blockchain instance with a genesis block.

Required:
  --InstanceId   Unique ID for this blockchain instance (auto-generated if omitted)

Example:
  novachain init --InstanceId 1001`,
		Run: func(cmd *cobra.Command, args []string) {
			if instanceID == "" {
				instanceID = fmt.Sprintf("%d", time.Now().Unix())
				log.Infof("No instance ID provided. Generated default: %s", instanceID)
			}
			cli, err := cli.UpdateInstance(instanceID, false, LogFile)
			if err != nil {
				log.Fatal(err)
			}
			cli.CreateBlockchain(chainData)
			log.Infof("✅ Blockchain initialized successfully (Instance ID: %s)", instanceID)
			cli.Blockchain.Database.Close()
		},
	}

	// -----------------------
	// WALLET COMMAND
	// -----------------------
	walletCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Manage wallets",
		Long:  `Create, list, or check wallet balances`,
	}

	walletCmd.AddCommand(
		&cobra.Command{
			Use:   "new",
			Short: "Create a new wallet",
			Run: func(cmd *cobra.Command, args []string) {
				cli.CreateWallet()
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all wallet addresses",
			Run: func(cmd *cobra.Command, args []string) {
				cli.ListWallet()
			},
		},
		&cobra.Command{
			Use:   "balance",
			Short: "Check wallet balance",
			Run: func(cmd *cobra.Command, args []string) {
				if address == "" {
					log.Fatal("--Address flag is required")
				}
				if instanceID == "" {
					log.Fatal("--InstanceId flag is required")
				}
				cli, err := cli.UpdateInstance(instanceID, true, LogFile)
				if err != nil {
					log.Fatal(err)
				}
				cli.GetBalance(address)
			},
		},
	)

	// -----------------------
	// NODE
	// -----------------------
	var (
		minerAddress string
		miner        bool
		fullNode     bool
		listenPort   string
	)

	nodeCmd := &cobra.Command{
		Use:   "startNode",
		Short: "Start a NovaChain node",
		Long: `Start a full node or mining node in the network.

Required:
  --Port        Port to listen on
  --InstanceId  Blockchain instance ID

Optional:
  --Miner       Run as mining node (true/false)
  --Address     Miner address (required if --Miner=true)
  --SeedPeer    Enable seed peer discovery
  --RPC         Enable JSON-RPC server

If no flags are provided, node runs as full node by default.`,
		Run: func(cmd *cobra.Command, args []string) {
			if instanceID == "" {
				log.Fatal("❌ You must run 'novachain init' before starting a node")
			}
			if listenPort == "" {
				log.Fatal("--Port flag is required")
			}
			if miner && minerAddress == "" {
				log.Fatal("Miner address is required when --Miner is true")
			}
			if !miner && !fullNode {
				fullNode = true
			}

			logPath := fmt.Sprintf("/logs/console_%s.log", instanceID)
			_ = utilLog.ClearLogFile(logPath)
			log.Infof("🔧 Starting node (Instance: %s, Port: %s, FullNode: %v, Miner: %v)", instanceID, listenPort, fullNode, miner)

			cli, err := cli.UpdateInstance(instanceID, false, LogFile)
			if err != nil {
				log.Fatal(err)
			}

			cli.StartNode(LogFile, chainData, listenPort, minerAddress, miner, fullNode, isSeedPeer, func(net *p2p.Network) {
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

	nodeCmd.Flags().StringVar(&listenPort, "Port", conf.ListenPort, "Node listening port")
	nodeCmd.Flags().StringVar(&minerAddress, "Address", conf.MinerAddress, "Miner address")
	nodeCmd.Flags().BoolVar(&miner, "Miner", conf.Miner, "Enable mining mode")
	nodeCmd.Flags().BoolVar(&fullNode, "Fullnode", conf.FullNode, "Run as full node")
	nodeCmd.Flags().BoolVar(&isSeedPeer, "SeedPeer", false, "Enable seed peer discovery")

	// -----------------------
	// ROOT COMMAND
	// -----------------------
	rootCmd := &cobra.Command{
		Use:   "novachain",
		Short: "NovaChain CLI",
		Long: `Command-line interface for running and interacting with the NovaChain blockchain.

Examples:
  1. Initialize blockchain:
     novachain init --InstanceId 1001

  2. Start a node (default full node):
     novachain startNode --Port 3000 --InstanceId 1001

  3. Start a miner node:
     novachain startNode --Port 3000 --InstanceId 1001 --Miner true --Address <wallet_address>

  4. Wallet management:
     novachain wallet new
     novachain wallet list
     novachain wallet balance --Address <wallet_address> --InstanceId 1001
`,
	}

	rootCmd.PersistentFlags().StringVar(&address, "Address", "", "Wallet address")
	rootCmd.PersistentFlags().StringVar(&instanceID, "InstanceId", "", "Blockchain instance ID")
	rootCmd.PersistentFlags().StringVar(&rpcPort, "RPC-Port", "9000", "RPC server port")
	rootCmd.PersistentFlags().StringVar(&rpcAddress, "RPC-Addr", "127.0.0.1", "RPC server address")
	rootCmd.PersistentFlags().BoolVar(&rpcEnabled, "RPC", false, "Enable JSON-RPC server")
	rootCmd.PersistentFlags().StringVar(&rpcMode, "RPC-Mode", "http", "RPC mode: http, tcp, both")
	rootCmd.PersistentFlags().StringVar(&chainData, "ChainData", "", "Chain data")
	rootCmd.PersistentFlags().StringVar(&LogFile, "LogFile", "", "Log data")

	rootCmd.AddCommand(initCmd, walletCmd, nodeCmd)

	if len(os.Args) == 1 {
		cui.Start(&cli, "config.json")
		return
	}

	rootCmd.Execute()
}
