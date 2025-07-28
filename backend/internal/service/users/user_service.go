package users

import (
	"errors"

	"github.com/miraklik/CargoLedger/internal/logger"
	"github.com/miraklik/CargoLedger/internal/models/users"
	"gorm.io/gorm"
)

var _ UsersInterface = (*UserService)(nil)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService { return &UserService{db: db} }

func (us *UserService) CreateUser(user *users.User) error {
	if err := us.db.Create(user).Error; err != nil {
		logger.Log.Errorf("Error creating user: %v", err)
		return err
	}

	logger.Log.Infof("User created successfully (id=%d)", user.ID)
	return nil
}

func (us *UserService) GetUser(id uint) (*users.User, error) {
	var user users.User

	err := us.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Warnf("User with ID %d not found", id)
		return nil, err
	}
	if err != nil {
		logger.Log.Errorf("Failed to get user (id=%d): %v", id, err)
		return nil, err
	}

	return &user, nil
}
