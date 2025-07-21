package cargo

import (
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type CargoStatus string

const (
	Created   CargoStatus = "created"
	InTransit CargoStatus = "inTransit"
	Delivered CargoStatus = "delivered"
	Cancelled CargoStatus = "cancelled"
)

type Cargo struct {
	gorm.Model
	Sender              common.Address `json:"sender"`
	Carrier             common.Address `json:"carrier"`
	Receiver            common.Address `json:"receiver"`
	DescriptionIpfsHash string         `json:"descriptionIpfsHash"`
	Status              CargoStatus    `json:"status"`
	Timestamp           int64          `json:"timestamp"`
}
