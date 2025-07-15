package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/miraklik/CargoLedger/internal/logger"
)

type Config struct {
	RPC struct {
		Rpc string
	}

	Database struct {
		Db_host string
		Db_user string
		Db_pass string
		Db_name string
		Db_port string
	}
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Errorf("Error loading .env file: %v", err)
		return nil, err
	}

	var cfg Config
	cfg.RPC.Rpc = os.Getenv("INFURA_API")

	cfg.Database.Db_host = os.Getenv("DB_HOST")
	cfg.Database.Db_user = os.Getenv("DB_USER")
	cfg.Database.Db_pass = os.Getenv("DB_PASS")
	cfg.Database.Db_name = os.Getenv("DB_NAME")
	cfg.Database.Db_port = os.Getenv("DB_PORT")

	return &cfg, nil
}
