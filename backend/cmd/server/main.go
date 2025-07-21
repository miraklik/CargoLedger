package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miraklik/CargoLedger/internal/handler"
	"github.com/miraklik/CargoLedger/internal/logger"
	"github.com/miraklik/CargoLedger/internal/service"
	"github.com/miraklik/CargoLedger/internal/storage"
)

func main() {
	db, err := storage.InitDB()
	if err != nil {
		logger.Log.Fatalf("Error initializing DB: %v", err)
	}

	services := service.NewCargoService(db)
	handlers := handler.NewCargoHandler(services)

	r := gin.Default()
	router := r.Group("/v1")
	router.POST("/cargo", handlers.CreateCargo)
	router.PUT("/cargo/:id", handlers.UpdateCargo)
	router.GET("/cargo/:id", handlers.GetCargoById)

	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatalf("Error starting server: %v", err)
	}
}
