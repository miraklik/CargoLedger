package users

import "gorm.io/gorm"

type UserType string

const (
	Sender   UserType = "Sender"
	Carrier  UserType = "Carrier"
	Receiver UserType = "Receiver"
)

type User struct {
	gorm.Model
	Address string   `json:"address" gorm:"uniqueIndex;not null"`
	Type    UserType `json:"type" gorm:"type:TypeUser;not null"`
}
