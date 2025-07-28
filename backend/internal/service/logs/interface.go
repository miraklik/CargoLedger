package logs

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LogsInterface interface {
	SaveEvent(eventType string, cargoID uint, txHash string, blockNumber uint64, index uint, sender string, rawData any) error
	ListenLogs(ctx context.Context, client *ethclient.Client, contractAddress string) error
}
