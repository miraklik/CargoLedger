package blockchain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/miraklik/CargoLedger/configs"
	"github.com/miraklik/CargoLedger/internal/logger"
)

var (
	cfg *configs.Config
)

func InitRPC() (*ethclient.Client, error) {
	client, err := ethclient.Dial(cfg.RPC.Rpc)
	if err != nil {
		logger.Log.Errorf("Error connecting to RPC: %v", err)
		return nil, err
	}

	logger.Log.Info("Connected to RPC")
	return client, nil
}
