package storage

import (
	"github.com/miraklik/CargoLedger/configs"
	"github.com/miraklik/CargoLedger/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	cfg, err := configs.Load()
	if err != nil {
		logger.Log.Fatalf("Error loading config: %v", err)
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.Db.Db_url), &gorm.Config{})
	if err != nil {
		logger.Log.Errorf("Error connecting to database: %v", err)
		return nil, err
	}

	logger.Log.Info("Connected to database")
	return db, nil
}
