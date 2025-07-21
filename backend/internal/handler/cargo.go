package handler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/miraklik/CargoLedger/internal/models/cargo"
	"github.com/miraklik/CargoLedger/internal/service"
	"net/http"
	"strconv"
)

type CargoHandler struct {
	CargoService service.CargoService
}

func NewCargoHandler(service service.CargoService) *CargoHandler {
	return &CargoHandler{
		CargoService: service,
	}
}

func (ch *CargoHandler) CreateCargo(c *gin.Context) {
	var req struct {
		Sender              common.Address    `json:"sender"`
		Carrier             common.Address    `json:"carrier"`
		Receiver            common.Address    `json:"receiver"`
		DescriptionIpfsHash string            `json:"description"`
		Status              cargo.CargoStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request body: " + err.Error()})
		return
	}

	newCargo := &cargo.Cargo{
		Sender:              req.Sender,
		Carrier:             req.Carrier,
		Receiver:            req.Receiver,
		DescriptionIpfsHash: req.DescriptionIpfsHash,
		Status:              req.Status,
	}

	if err := ch.CargoService.CreateCargo(newCargo); err != nil {
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
		Status cargo.CargoStatus `json:"status"`
	}

	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id: " + id})
		return
	}
	uuid := uint(uid)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind request body: " + err.Error()})
		return
	}

	if err := ch.CargoService.UpdateStatus(uuid, req.Status); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated cargo"})
}

func (ch *CargoHandler) GetCargoById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id: " + idParam})
		return
	}

	cargos, err := ch.CargoService.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cargo": cargos,
	})
}
