package cargo

import (
	"errors"
	"fmt"
	"github.com/miraklik/CargoLedger/internal/logger"
	"github.com/miraklik/CargoLedger/internal/models/cargo"
	"gorm.io/gorm"
)

var _ CargoInterface = (*CargoService)(nil)

type CargoService struct {
	db *gorm.DB
}

func NewCargoService(db *gorm.DB) *CargoService {
	return &CargoService{db: db}
}

func (cs *CargoService) CreateCargo(cargo *cargo.Cargo) error {
	if err := cs.db.Create(cargo).Error; err != nil {
		logger.Log.Errorf("Error creating cargo: %v", err)
		return err
	}

	logger.Log.Info("Cargo created successfully")
	return nil
}

func (cs *CargoService) GetById(id uint) (*cargo.Cargo, error) {
	var cg cargo.Cargo
	err := cs.db.First(&cg, id).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		logger.Log.Warnf("Cargo with ID %d not found", id)
		return nil, err
	case err != nil:
		logger.Log.Errorf("Failed to retrieve cargo with ID %d: %v", id, err)
		return nil, fmt.Errorf("get cargo: %w", err)
	}

	logger.Log.Infof("Cargo with ID %d retrieved successfully", id)
	return &cg, nil
}

func (cs *CargoService) UpdateStatus(id uint, status cargo.CargoStatus) (*cargo.Cargo, error) {
	var cg cargo.Cargo
	result := cs.db.Model(&cg).Where("id = ?", id).Update("status", status)
	if err := result.Error; err != nil {
		logger.Log.Errorf("Failed to update cargo status (id=%d): %v", id, err)
		return nil, fmt.Errorf("update status: %w", err)
	}

	if result.RowsAffected == 0 {
		logger.Log.Warnf("No cargo found with ID %d to update", id)
		return nil, gorm.ErrRecordNotFound
	}

	logger.Log.Infof("Cargo with ID %d status updated to %s", id, status)
	return &cg, nil
}
