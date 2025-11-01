package env

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	_, file, _, _ = runtime.Caller(0)

	root = filepath.Join(filepath.Dir(file), "../../../")
	once sync.Once
)

type Env struct {
	ServerPort                      string
	AppEnv                          string
	MaxTimeLog                      int64 // days
	Fullnode_RPC_URL                string
	DB_URL                          string
	Wallet_Signature_Expiry_Minutes int64 // minutes
	Jwt_Secret_Key                  string
	Jwt_TTL_Minutes                 int64
	Domain_Client                   string
	CheckSumLength                  int64
	Version                         int64
	Encode_data_secret_Key          []byte
	Sync_block_batch_size           int64
}

var Cfg *Env

func InitEnv() {
	once.Do(func() {
		Cfg = &Env{
			ServerPort:                      GetEnvAsString("SERVER_PORT", "3001"),
			AppEnv:                          GetEnvAsString("APP_ENV", "development"),
			MaxTimeLog:                      GetEnvAsInt("MAX_TIME_LOG", 30),
			Fullnode_RPC_URL:                GetEnvAsString("FULLNODE_RPC_URL", "http://0.0.0.0:9050/__jsonrpc"),
			DB_URL:                          GetEnvAsString("DATABASE_URL", ""),
			Wallet_Signature_Expiry_Minutes: GetEnvAsInt("WALLET_SIGNATURE_EXPIRY_MINUTES", 1),
			Jwt_Secret_Key:                  GetEnvAsString("JWT_SECRET_KEY", "default_secret_key"),
			Jwt_TTL_Minutes:                 GetEnvAsInt("JWT_TTL_MINUTES", 1800),
			Domain_Client:                   GetEnvAsString("DOMAIN_CLIENT", "localhost"),
			CheckSumLength:                  GetEnvAsInt("CHECK_SUM_LENGTH", 4),
			Version:                         GetEnvAsInt("VERSION", 0x00),
			Encode_data_secret_Key:          GetEnvAsBytes("ENCODE_DATA_SECRET_KEY", []byte("")),
			Sync_block_batch_size:           GetEnvAsInt("SYNC_BLOCK_BATCH_SIZE", 100),
		}
	})
}

func GetEnvVariable(key string) string {
	envPath := filepath.Join(root, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Panic("⚠️ No .env file found: " + err.Error())
	}

	return os.Getenv(key)
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int64) int64 {
	valueStr := GetEnvVariable(key)

	if value, err := strconv.Atoi(valueStr); err == nil {
		return int64(value)
	}

	return defaultValue
}

func GetEnvAsString(key, defaultValue string) string {
	valueStr := GetEnvVariable(key)

	if valueStr != "" {
		return valueStr
	}

	return defaultValue
}

func GetEnvAsBool(key string, defaultValue bool) bool {
	valueStr := GetEnvVariable(key)

	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultValue
}

func GetEnvAsBytes(key string, defaultValue []byte) []byte {
	valueStr := GetEnvVariable(key)
	if valueStr != "" {
		return []byte(valueStr)
	}

	return defaultValue
}
