package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/miraklik/CargoLedger/configs"
	"github.com/miraklik/CargoLedger/internal/blockchain"
	"github.com/miraklik/CargoLedger/internal/handler"
	"github.com/miraklik/CargoLedger/internal/logger"
	"github.com/miraklik/CargoLedger/internal/service"
	"github.com/miraklik/CargoLedger/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		logger.Log.Fatalf("Error loading config: %v", err)
	}

	db, err := storage.InitDB()
	if err != nil {
		logger.Log.Fatalf("Error initializing DB: %v", err)
	}

	client, err := blockchain.InitRPC()
	if err != nil {
		logger.Log.Fatalf("Error initializing RPC: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cargoServices := service.NewCargoService(db)
	cargoHandlers := handler.NewCargoHandler(cargoServices)
	logs := service.NewLogService(db)

	contractAddr := string(cfg.Contract.Contract_Address)

	go func() {
		if err := logs.ListenLogs(ctx, client, contractAddr); err != nil {
			logger.Log.Errorf("Log listener stopped with error: %v", err)
		}
	}()

	r := gin.Default()
	router := r.Group("/v1")
	router.POST("/cargo", cargoHandlers.CreateCargo)
	router.PUT("/cargo/:id", cargoHandlers.UpdateCargo)
	router.GET("/cargo/:id", cargoHandlers.GetCargoById)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		logger.Log.Info("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Log.Info("Server exited")
}
