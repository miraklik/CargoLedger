package cargo

import "github.com/miraklik/CargoLedger/internal/models/cargo"

type CargoInterface interface {
	CreateCargo(cargo *cargo.Cargo) error
	GetById(id uint) (*cargo.Cargo, error)
	UpdateStatus(id uint, status cargo.CargoStatus) (*cargo.Cargo, error)
}
