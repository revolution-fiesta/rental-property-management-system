package controllers

import (
    "net/http"
    "rental-property-management-system/internal/database"
    "rental-property-management-system/internal/models"
    "github.com/gin-gonic/gin"
)

// 查询房间分配关系，包含房间信息
func GetAllRelationships(c *gin.Context) {
    var relationships []models.Relationship

    // 使用 Preload 来加载关联的房间信息
    if err := database.DB.Preload("Room").Find(&relationships).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve relationships"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "relationships": relationships,
    })
}
