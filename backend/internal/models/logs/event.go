package logs

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	EventType     string `gorm:"size:50"`
	CargoID       uint
	TxHash        string `gorm:"size:255"`
	BlockNumber   uint64
	EventIndex    uint
	SenderAddress string         `gorm:"size:50"`
	Data          datatypes.JSON `gorm:"type:jsonb"`
}
