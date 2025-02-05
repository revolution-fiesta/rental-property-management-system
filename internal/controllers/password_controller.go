package controllers

import (
	"math/rand"
	"net/http"
	"time"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"
	"github.com/gin-gonic/gin"
)

func GenerateTempPassword(c *gin.Context) {
	var request struct {
		RoomID uint `json:"room_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	password := models.Password{
		RoomID:    request.RoomID,
		Password:  generateRandomPassword(6),
		IsTemp:    true,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := database.DB.Create(&password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password":    password.Password,
		"expires_at": password.ExpiresAt,
	})
}

func generateRandomPassword(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}