package controller

import (
	"net/http"

	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

// 根据用户 ID 查询房间分配关系
func ListRelationships(c *gin.Context) {

	var relationships []store.Relationship

	// 使用 Preload 来加载关联的房间信息
	if err := store.GetDB().Preload("Room").Find(&relationships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve relationships"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"relationships": relationships,
	})
}
