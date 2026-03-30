package handlers

import (
	"net/http"

	"serverpoint/config"
	"serverpoint/database"
	"serverpoint/models"
	"serverpoint/services"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	Config *config.Config
}

type createServerRequest struct {
	Name       string `json:"name" binding:"required"`
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port"`
	Username   string `json:"username" binding:"required"`
	AuthType   string `json:"auth_type" binding:"required,oneof=password key"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
}

func (h *ServerHandler) List(c *gin.Context) {
	var servers []models.Server
	database.DB.Order("created_at desc").Find(&servers)
	c.JSON(http.StatusOK, servers)
}

func (h *ServerHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var server models.Server
	if err := database.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}
	c.JSON(http.StatusOK, server)
}

func (h *ServerHandler) Create(c *gin.Context) {
	var req createServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Port == 0 {
		req.Port = 22
	}

	userID, _ := c.Get("user_id")

	server := models.Server{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		AuthType: req.AuthType,
		UserID:   userID.(uint),
	}

	if req.AuthType == "password" && req.Password != "" {
		encrypted, err := services.Encrypt(req.Password, h.Config.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt credentials"})
			return
		}
		server.EncryptedPassword = encrypted
	}

	if req.AuthType == "key" && req.PrivateKey != "" {
		encrypted, err := services.Encrypt(req.PrivateKey, h.Config.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt credentials"})
			return
		}
		server.EncryptedKey = encrypted
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create server"})
		return
	}

	c.JSON(http.StatusCreated, server)
}

func (h *ServerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var server models.Server
	if err := database.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	// Clean up related data
	database.DB.Where("server_id = ?", server.ID).Delete(&models.Metric{})
	database.DB.Where("server_id = ?", server.ID).Delete(&models.Process{})
	database.DB.Delete(&server)

	c.JSON(http.StatusOK, gin.H{"message": "server deleted"})
}
