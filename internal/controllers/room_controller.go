package controllers

import (
	"net/http"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAvailableRooms(c *gin.Context) {
	var rooms []models.Room

	result := database.DB.Where("is_deleted = ?", false).Find(&rooms)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rooms})
}