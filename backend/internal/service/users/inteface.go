package users

import "github.com/miraklik/CargoLedger/internal/models/users"

type UsersInterface interface {
	CreateUser(user *users.User) error
	GetUser(id uint) (*users.User, error)
}
