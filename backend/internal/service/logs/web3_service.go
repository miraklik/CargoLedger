package logs

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/miraklik/CargoLedger/internal/logger"
	"github.com/miraklik/CargoLedger/internal/models/logs"
	"gorm.io/gorm"
	"math/big"
	"os"
	"strings"
)

type LogService struct {
	db *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{db: db}
}

func (s *LogService) SaveEvent(eventType string, cargoID uint, txHash string, blockNumber uint64, index uint, sender string, rawData any) error {
	dataJSON, err := json.Marshal(rawData)
	if err != nil {
		return err
	}

	event := logs.Event{
		EventType:     eventType,
		CargoID:       cargoID,
		TxHash:        txHash,
		BlockNumber:   blockNumber,
		EventIndex:    index,
		SenderAddress: sender,
		Data:          dataJSON,
	}

	return s.db.Create(&event).Error
}

func (s *LogService) ListenLogs(ctx context.Context, client *ethclient.Client, contractAddress string) error {
	abiJSON, err := os.ReadFile("./blockchain/abi/CargoLedger_ABI.json")
	if err != nil {
		logger.Log.Printf("Failed to read abi: %v", err)
		return err
	}

	abiJSONString := string(abiJSON)
	contractAbi, err := abi.JSON(strings.NewReader(abiJSONString))
	if err != nil {
		logger.Log.Printf("Failed to parse contract ABI: %v", err)
		return err
	}

	address := common.HexToAddress(contractAddress)

	logsCh := make(chan types.Log)
	query := ethereum.FilterQuery{Addresses: []common.Address{address}}

	sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		logger.Log.Printf("Failed to subscribe to logs: %v", err)
		return err
	}

	logger.Log.Info("🟢 Event listener started")

	for {
		select {
		case err := <-sub.Err():
			logger.Log.Printf("Log subscription error: %v", err)
			return err
		case vLog := <-logsCh:
			eventID := vLog.Topics[0]
			var eventName string

			for name, evt := range contractAbi.Events {
				if evt.ID == eventID {
					eventName = name
					break
				}
			}

			if eventName == "" {
				logger.Log.Warn("Unknown event received")
				continue
			}

			logger.Log.Infof("📦 Event received: %s", eventName)

			switch eventName {
			case "CargoCreated":
				var parsed struct {
					CargoID  *big.Int
					Sender   common.Address
					Receiver common.Address
					IpfsHash string
				}

				err := contractAbi.UnpackIntoInterface(&parsed, eventName, vLog.Data)
				if err != nil {
					logger.Log.Errorf("Failed to unpack %s: %v", eventName, err)
					continue
				}

				err = s.SaveEvent(eventName,
					uint(parsed.CargoID.Uint64()),
					vLog.TxHash.Hex(),
					vLog.BlockNumber,
					uint(vLog.Index),
					vLog.Address.Hex(),
					parsed,
				)

			case "CargoUpdated":
				var parsed struct {
					CargoID *big.Int
					NewHash string
					Updater common.Address
				}

				err := contractAbi.UnpackIntoInterface(&parsed, eventName, vLog.Data)
				if err != nil {
					logger.Log.Errorf("Failed to unpack %s: %v", eventName, err)
					continue
				}

				err = s.SaveEvent(eventName,
					uint(parsed.CargoID.Uint64()),
					vLog.TxHash.Hex(),
					vLog.BlockNumber,
					uint(vLog.Index),
					vLog.Address.Hex(),
					parsed,
				)
			default:
				logger.Log.Warnf("Unhandled event: %s", eventName)
				continue
			}

			if err != nil {
				logger.Log.Errorf("Failed to save event %s: %v", eventName, err)
			}
		}
	}
}
