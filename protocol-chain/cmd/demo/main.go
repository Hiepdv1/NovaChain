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

func main() {
	defer os.Exit(0)

	var conf = env.New()
	var address string
	var InstanceId string
	var rpcPort string
	var rpcAddress string
	var rpcMode string
	var rpc bool
	var isSeedPeer bool

	cli := utilCmd.CommandLine{
		Blockchain: &blockchain.Blockchain{
			Database:   nil,
			InstanceId: InstanceId,
		},
		P2P: nil,
	}

	/*
	* INIT COMMAND
	 */
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize blockchain and create the genesis block",
		Long: `Initialize the blockchain instance with a genesis block.

Required:
  --InstanceId   Unique numeric ID for the blockchain instance (auto-generated if not provided).

Example:
  blockchain init --InstanceId <your_id>
`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if InstanceId == "" {
				InstanceId = fmt.Sprintf("%d", time.Now().Unix())
				log.Infof("No instance ID provided, Generated default ID: %s", InstanceId)
			}
			cli, err := cli.UpdateInstance(InstanceId, false)
			if err != nil {
				log.Error(err)
			}
			cli.CreateBlockchain()
			log.Infof("✅ Blockchain initialized successfully with instance ID: %s", InstanceId)
			cli.Blockchain.Database.Close()
		},
	}

	/*
	* WALLET COMMAND
	 */

	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Wallet management commands",
		Long: `Use this command group to manage wallets.

Subcommands:
  - new         Create a new wallet.
  - listAddress List all available wallet addresses.
  - balance     Get the balance of a specific wallet address.
`,
	}

	var newWalletCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new wallet",
		Long:  `Generate a new wallet and print its address.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.CreateWallet()
		},
	}

	var listWalletAddressCmd = &cobra.Command{
		Use:   "list",
		Short: "List all wallet addresses",
		Long:  `List all wallet addresses stored locally.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli.ListWallet()
		},
	}

	var walletBalanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Check balance of a wallet address",
		Long: `Get the current balance of a given wallet address.

Required:
  --Address   The wallet address to check balance (must be 34 characters).

Example:
  blockchain wallet balance --Address <wallet_address>
`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := cli.UpdateInstance(InstanceId, true)
			if err != nil {
				log.Error(err)
			}
			cli.GetBalance(address)
		},
	}

	walletCmd.AddCommand(
		newWalletCmd,
		listWalletAddressCmd,
		walletBalanceCmd,
	)

	/*
	* UTXOS COMMAND
	 */
	var computeUtxosCmd = &cobra.Command{
		Use:   "computeutxos",
		Short: "Recalculate unspent transaction outputs (UTXOs)",
		Long: `Rebuild the entire set of unspent transaction outputs (UTXOs) from the blockchain.

Should be run when UTXOs are out-of-sync or after data corruption recovery.`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := cli.UpdateInstance(InstanceId, false)
			if err != nil {
				log.Error(err)
			}
			cli.ComputeUTXOs()
		},
	}

	/*
	* PRINT COMMAND
	 */
	var printCmd = &cobra.Command{
		Use:   "print",
		Short: "Print all blocks in the blockchain",
		Long:  `Display all blocks stored in the blockchain, including transactions and metadata.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if InstanceId == "" {
				fmt.Print("Instance Required")
			}

			cli, err := cli.UpdateInstance(InstanceId, true)
			if err != nil {
				log.Error(err)
			}
			cli.PrintBlockchain()
		},
	}

	var utxosCmd = &cobra.Command{
		Use:  "utxos",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cli, err := cli.UpdateInstance(InstanceId, true)
			if err != nil {
				log.Error(err)
			}
			cli.PrintUtxos()
		},
	}

	var bestChainCmd = &cobra.Command{
		Use:  "best-chain",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if InstanceId == "" {
				fmt.Print("Instance Required")
			}

			cli, err := cli.UpdateInstance(InstanceId, true)
			if err != nil {
				log.Error(err)
			}
			cli.PrintBestChain()
		},
	}

	/*
	* NODE COMMAND
	 */
	var minerAddress string
	var miner bool
	var fullNode bool
	var listenPort string
	var nodeCmd = &cobra.Command{
		Use:   "startNode",
		Short: "Start a node in the blockchain network",
		Long: `Start a full node or mining node and connect to the P2P network.

Required:
  --Port        Port to listen on.

  --InstanceId  Unique numeric ID for the blockchain instance
  
  Optional:
  --Address     Miner address (if --miner is set, must be a valid wallet address of 34 characters).
  
  --Miner       Set this node as a miner (true/false).

  --Fullnode    Set this node as a full node (true/false).

  --RPC         Enable HTTP JSON-RPC server.

  --RPC-Port    Port for JSON-RPC server (default 9000).

  --RPC-Addr    Address interface for RPC server (default localhost) (example: 127.0.0.1, 0.0.0.0, ...v.v).
  
  --RPC-Mode    Set the JSON-RPC mode: "http", "tcp", or "both".
                Use "http" to expose RPC over HTTP (e.g., via browser or REST-like requests),
                "tcp" for lower-level socket communication (e.g., CLI clients), 
                or "both" to enable both interfaces.
Example:
  blockchain startNode --Port 3000 --InstanceId 3000 --miner true --address <wallet_address>
`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			logPath := fmt.Sprintf("/logs/console_%s.log", InstanceId)
			_ = utilLog.ClearLogFile(logPath)

			cli, err := cli.UpdateInstance(InstanceId, false)

			if err != nil {
				log.Error(err)
			}

			if miner && len(minerAddress) == 0 {
				log.Fatal("Miner address is required when starting a miner node")
			}

			log.Infof("🔧 Starting node with instanceId: %s on port: %s", InstanceId, listenPort)

			cli.StartNode(listenPort, minerAddress, miner, fullNode, isSeedPeer, func(net *p2p.Network) {
				log.Info("✅ Node started successfully and joined the network")
				if miner {
					log.Info("Miner mode")
				}
				if rpc {
					cli.P2P = net
					go jsonrpc.StartServer(cli, rpc, rpcPort, rpcAddress, rpcMode)
				}
			})
		},
	}

	nodeCmd.Flags().StringVar(&listenPort, "Port", conf.ListenPort, "Node listening port")
	nodeCmd.Flags().StringVar(&minerAddress, "Address", conf.MinerAddress, "Set miner address")
	nodeCmd.Flags().BoolVar(&miner, "Miner", conf.Miner, "Set as true if you are joining the network as a miner")
	nodeCmd.Flags().BoolVar(&fullNode, "Fullnode", conf.FullNode, "Set as true if you are joining the network as a miner")
	nodeCmd.Flags().BoolVar(&isSeedPeer, "SeedPeer", false, "Enable the Seed Peer")
	/*
	* SEND COMMAND
	 */

	var mineNow bool
	var sendFrom, sendTo string
	var amount, fee float64

	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send tokens from one wallet address to another",
		Long: `Send a transaction from one wallet to another.

Required Flags:
  --SendFrom    Sender's wallet address (must be 34 characters).

  --SendTo      Recipient's wallet address (must be 34 characters).

  --Amount      Amount of token to send (must be > 0).

  --Fee         Transaction fee (default is 1/PER_COIN).

Optional Flag:
  --Mine        If set to true, your node will mine the transaction immediately.

Example:
  blockchain send --SendFrom <sender_address> --SendTo <receiver_address> --Amount 10 --Fee 0.01 --Mine
`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// cli.Send(sendFrom, sendTo, amount, fee, mineNow)
		},
	}

	sendCmd.Flags().StringVar(&sendFrom, "SendFrom", "", "Sender's wallet address")
	sendCmd.Flags().StringVar(&sendTo, "SendTo", "", "Receiver's wallet address")
	sendCmd.Flags().Float64Var(&amount, "Amount", float64(0), "Amount of token to send")
	sendCmd.Flags().Float64Var(&fee, "Fee", float64(1/blockchain.PER_COIN), "Fee for this transaction")
	sendCmd.Flags().BoolVar(&mineNow, "Mine", false, "Set if you want your node to mine the transaction instantly")

	var rootCmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Blockchain CLI tool",
		Long: `Blockchain CLI
A command-line interface for managing and interacting with the blockchain system.

Global Required Flags:
  --Address      Your wallet address (must be 34 characters).

  --InstanceId   Blockchain instance ID (a numeric string to isolate data per instance).

Example Usage:

  1. Initialize blockchain:
     blockchain init --address <wallet_address> --InstanceId <your_id>

  2. Start node:
     blockchain startNode --port 3000 --address <wallet_address> --miner true

  3. Create a new wallet:
     blockchain wallet new

  4. Send tokens:
     blockchain send --SendFrom <sender_address> --SendTo <receiver_address> --Amount 5 --Fee 0.01 --Mine

  5. View blockchain:
     blockchain print
`,
	}

	rootCmd.PersistentFlags().StringVar(&address, "Address", "", "Your wallet address")
	rootCmd.PersistentFlags().StringVar(&InstanceId, "InstanceId", "", "Blockchain instance")

	/*
	* HTTP FLAGS
	 */
	rootCmd.PersistentFlags().StringVar(&rpcPort, "RPC-Port", "9000", "HTTP-RPC server listening port (default: 9000)")
	rootCmd.PersistentFlags().StringVar(&rpcAddress, "RPC-Addr", "127.0.0.1", "HTTP-RPC server listening interface (default: localhost)")
	rootCmd.PersistentFlags().BoolVar(&rpc, "RPC", false, "Enable the HTTP-RPC server")
	rootCmd.PersistentFlags().StringVar(&rpcMode, "RPC-Mode", "http", "Select JSON-RPC mode: 'http', 'tcp', or 'both' (default: 'http')")

	rootCmd.AddCommand(
		initCmd,
		walletCmd,
		computeUtxosCmd,
		utxosCmd,
		sendCmd,
		printCmd,
		nodeCmd,
		bestChainCmd,
	)

	rootCmd.Execute()

}
