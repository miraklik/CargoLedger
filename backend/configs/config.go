package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/miraklik/CargoLedger/internal/logger"
)

type Config struct {
	Db       Database
	Rpc      RPC
	Contract Contract
}

type RPC struct {
	Rpc string
}

type Database struct {
	Db_url string
}

type Contract struct {
	Contract_Address string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Db: Database{
			Db_url: getEnv("DB_URL", ""),
		},
		Rpc: RPC{
			Rpc: getEnv("RPC_URL", ""),
		},
		Contract: Contract{
			Contract_Address: getEnv("CONTRACT_ADDR", ""),
		},
	}

	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		if defaultVal == "" {
			logger.Log.Fatalf("Failed to get %s from environment", key)
		}
		return defaultVal
	}
	return val
}
