package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/miraklik/CargoLedger/configs"
	"github.com/miraklik/CargoLedger/internal/blockchain"
	cargoHandler "github.com/miraklik/CargoLedger/internal/handler/cargo"
	usersHandler "github.com/miraklik/CargoLedger/internal/handler/users"
	"github.com/miraklik/CargoLedger/internal/logger"
	cargoService "github.com/miraklik/CargoLedger/internal/service/cargo"
	"github.com/miraklik/CargoLedger/internal/service/logs"
	usersService "github.com/miraklik/CargoLedger/internal/service/users"
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

	cargoServices := cargoService.NewCargoService(db)
	cargoHandlers := cargoHandler.NewCargoHandler(*cargoServices)

	userService := usersService.NewUserService(db)
	userHandler := usersHandler.NewUserHandler(*userService)
	loggers := logs.NewLogService(db)

	contractAddr := string(cfg.Contract.Contract_Address)

	go func() {
		if err := loggers.ListenLogs(ctx, client, contractAddr); err != nil {
			logger.Log.Errorf("Log listener stopped with error: %v", err)
			return
		}
	}()

	r := gin.Default()
	cargoGroup := r.Group("/cargo")
	{
		cargoGroup.POST("/", cargoHandlers.CreateCargo)
		cargoGroup.PUT("/:id", cargoHandlers.UpdateCargo)
		cargoGroup.GET("/:id", cargoHandlers.GetCargoById)
	}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUser)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
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
