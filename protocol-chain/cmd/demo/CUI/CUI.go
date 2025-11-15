package cui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	utilCmd "core-blockchain/cmd/utils"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	InstanceID   string
	Port         string
	MinerAddress string
	SeedPeer     bool
	ChainData    string
	LogFile      string
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	return cfg, err
}

func Start(cli *utilCmd.CommandLine, configPath string) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Warnf("❌ Failed to load config, using defaults: %v", err)
		cfg = &Config{Port: "3000"}
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = fmt.Sprintf("%d", time.Now().Unix())
	}

	log.Infof("Generated InstanceID: %s", cfg.InstanceID)

	c2, err := cli.UpdateInstance(cfg.InstanceID, true, cfg.LogFile)
	if err != nil {
		log.Fatal(err)
	}

	c2.CreateBlockchain(cfg.ChainData)

	writeConfig(configPath, cfg)

	for {
		fmt.Println("======================================")
		fmt.Println("       🚀 NovaChain Command Center")
		fmt.Println("======================================")
		fmt.Println("1. Start full node")
		fmt.Println("2. Start miner node")
		fmt.Println("3. Create new wallet")
		fmt.Println("4. List wallets")
		fmt.Println("5. Exit")
		fmt.Print("Select option: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			startNodeProcess(cfg.InstanceID, cfg.Port, "", false, cfg.SeedPeer, cfg.ChainData, cfg.LogFile)

		case 2:
			if cfg.MinerAddress == "" {
				fmt.Println("❌ Missing MinerAddress in config.")
				break
			}
			startNodeProcess(cfg.InstanceID, cfg.Port, cfg.MinerAddress, true, cfg.SeedPeer, cfg.ChainData, cfg.LogFile)

		case 3:
			cli.CreateWallet()

		case 4:
			cli.ListWallet()

		case 5:
			fmt.Println("👋 Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func clearTerminal() {
	cmd := exec.Command("cmd", "/c", "cls")
	// cmd := exec.Command("clear") // Linux/Mac
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func startNodeProcess(instanceID, port, minerAddr string, isMiner, seed bool, chainData, logFile string) {
	clearTerminal()
	args := []string{
		"startNode",
		"--InstanceId", instanceID,
		"--Port", port,
		"--ChainData", chainData,
		"--LogFile", logFile,
	}

	if isMiner {
		args = append(args, "--Miner", "true", "--Address", minerAddr)
	}
	if seed {
		args = append(args, "--Seed", "true")
	}

	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Infof("🚀 Node started in a separate process (PID: %d)", cmd.Process.Pid)
}

func writeConfig(path string, cfg *Config) {
	f, _ := os.Create(path)
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.Encode(cfg)
	f.Close()
}
