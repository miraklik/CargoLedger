package service

import (
	"github.com/miraklik/CargoLedger/internal/logger"
	"github.com/miraklik/CargoLedger/internal/models/cargo"
	"gorm.io/gorm"
)

type CargoService struct {
	db *gorm.DB
}

func NewCargoService(db *gorm.DB) *CargoService {
	return &CargoService{db: db}
}

func (c *CargoService) CreateCargo(cargo *cargo.Cargo) error {
	if err := c.db.Create(cargo).Error; err != nil {
		logger.Log.Errorf("Error creating cargo: %v", err)
		return err
	}

	logger.Log.Info("Cargo created successfully")
	return nil
}

func (c *CargoService) GetCargoById(id uint) (*cargo.Cargo, error) {
	var cargos cargo.Cargo
	if err := c.db.First(&cargos, id).Error; err != nil {
		logger.Log.Errorf("Error getting cargo: %v", err)
		return nil, err
	}

	logger.Log.Infof("Cargo with id %d found successfully", id)
	return &cargos, nil
}

func (c *CargoService) CargoUpdateStatus(id uint, status cargo.CargoStatus) error {
	var cargos cargo.Cargo
	if err := c.db.Model(&cargos).Where("id = ?", id).Update("status", status).Error; err != nil {
		logger.Log.Errorf("Error updating cargo status: %v", err)
		return err
	}

	logger.Log.Infof("Cargo with id %d updated successfully", id)
	return nil
}
