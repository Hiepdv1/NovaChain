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

	root = filepath.Join(filepath.Dir(file), "../../../")
)

type Env struct {
	ServerPort       string
	AppEnv           string
	MaxTimeLog       int64 // days
	Fullnode_RPC_URL string
	DB_URL           string
}

func New() *Env {
	return &Env{
		ServerPort:       GetEnvAsString("SERVER_PORT", "3001"),
		AppEnv:           GetEnvAsString("APP_ENV", "development"),
		MaxTimeLog:       GetEnvAsInt("MAX_TIME_LOG", 30),
		Fullnode_RPC_URL: GetEnvAsString("FULLNODE_RPC_URL", "http://0.0.0.0:9050/__jsonrpc"),
		DB_URL:           GetEnvAsString("DATABASE_URL", ""),
	}
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
