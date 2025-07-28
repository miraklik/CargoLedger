package cargo

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/miraklik/CargoLedger/internal/logger"
	cargoModel "github.com/miraklik/CargoLedger/internal/models/cargo"
	cargoService "github.com/miraklik/CargoLedger/internal/service/cargo"
	"net/http"
	"strconv"
)

type CargoHandler struct {
	CargoService cargoService.CargoService
}

func NewCargoHandler(service cargoService.CargoService) *CargoHandler {
	return &CargoHandler{
		CargoService: service,
	}
}

func (ch *CargoHandler) CreateCargo(c *gin.Context) {
	var req struct {
		Sender              common.Address         `json:"sender"`
		Carrier             common.Address         `json:"carrier"`
		Receiver            common.Address         `json:"receiver"`
		DescriptionIpfsHash string                 `json:"description"`
		Status              cargoModel.CargoStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Failed to bind JSON body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request body: " + err.Error()})
		return
	}

	newCargo := &cargoModel.Cargo{
		Sender:              req.Sender,
		Carrier:             req.Carrier,
		Receiver:            req.Receiver,
		DescriptionIpfsHash: req.DescriptionIpfsHash,
		Status:              req.Status,
	}

	if err := ch.CargoService.CreateCargo(newCargo); err != nil {
		logger.Log.Errorf("Failed to create cargo: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created cargo",
		"cargo":   req,
	})
}

func (ch *CargoHandler) UpdateCargo(c *gin.Context) {
	var req struct {
		Status cargoModel.CargoStatus `json:"status"`
	}

	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logger.Log.Errorf("Failed to parse id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id: " + id})
		return
	}
	uuid := uint(uid)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Failed to bind JSON body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request body: " + err.Error()})
		return
	}

	cargos, err := ch.CargoService.UpdateStatus(uuid, req.Status)
	if err != nil {
		logger.Log.Errorf("Failed to update status: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated cargo",
		"cargo":   cargos,
	})
}

func (ch *CargoHandler) GetCargoById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logger.Log.Errorf("Failed to parse id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id: " + idParam})
		return
	}

	cargos, err := ch.CargoService.GetById(uint(id))
	if err != nil {
		logger.Log.Errorf("Failed to get cargo: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cargo": cargos,
	})
}
