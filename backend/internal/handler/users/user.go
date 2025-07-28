package users

import (
	"github.com/gin-gonic/gin"
	"github.com/miraklik/CargoLedger/internal/logger"
	usersModel "github.com/miraklik/CargoLedger/internal/models/users"
	usersService "github.com/miraklik/CargoLedger/internal/service/users"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserService usersService.UserService
}

func NewUserHandler(userService usersService.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Address string              `json:"address"`
		Type    usersModel.UserType `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Errorf("Failed to bind JSON body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &usersModel.User{
		Address: req.Address,
		Type:    req.Type,
	}

	if err := uh.UserService.CreateUser(user); err != nil {
		logger.Log.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logger.Log.Warnf("Failed to parse user id %s", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.UserService.GetUser(uint(uid))
	if err != nil {
		logger.Log.Errorf("Failed to get user with id %d", uid)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Debugf("User with id %d retrieved successfully", uid)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
