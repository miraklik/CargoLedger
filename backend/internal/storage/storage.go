package storage

import (
	"fmt"

	"github.com/miraklik/CargoLedger/configs"
	"github.com/miraklik/CargoLedger/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	cfg, err := configs.Load()
	if err != nil {
		logger.Log.Errorf("Error loading config: %v", err)
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Database.Db_host, cfg.Database.Db_user, cfg.Database.Db_pass, cfg.Database.Db_name, cfg.Database.Db_port)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		logger.Log.Errorf("Error connecting to database: %v", err)
		return nil, err
	}

	logger.Log.Info("Connected to database")
	return db, nil
}
