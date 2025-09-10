package env

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	_, file, _, _ = runtime.Caller(0)

	root = filepath.Join(filepath.Dir(file), "../..")
)

type Config struct {
	WalletAddressCheckSum int64
	MinerAddress          string
	ListenPort            string
	Miner                 bool
	FullNode              bool
	SystemKey             string
	Wallet_Padding        string
}

func New() *Config {
	return &Config{
		WalletAddressCheckSum: GetEnvAsInt("WALLET_ADDRESS_CHECKSUM", 1),
		MinerAddress:          GetEnvAsStr("MINER_ADDRESS", ""),
		ListenPort:            GetEnvAsStr("LISTEN_PORT", ""),
		Miner:                 GetEnvAsBool("MINER", false),
		FullNode:              GetEnvAsBool("FULL_NODE", false),
		SystemKey:             GetEnvAsStr("SYSTEM_KEY", ""),
		Wallet_Padding:        GetEnvAsStr("WALLET_PADDING", ""),
	}
}

func GetEnvVariable(key string) string {

	envPath := filepath.Join(root, ".env")

	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetEnvAsInt(name string, defaultValue int64) int64 {
	valueStr := GetEnvVariable(name)

	if value, err := strconv.Atoi(valueStr); err == nil {
		return int64(value)
	}

	return defaultValue
}

func GetEnvAsStr(name string, defaultValue string) string {
	valueStr := GetEnvVariable(name)

	if valueStr != "" {
		return valueStr
	}

	return defaultValue
}

func GetEnvAsBool(name string, defaultValue bool) bool {
	valueStr := GetEnvVariable(name)

	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultValue
}
