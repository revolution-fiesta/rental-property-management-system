package controller

import (
	"net/http"

	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

// 根据用户 ID 查询房间分配关系
func ListRelationships(c *gin.Context) {
	var relationships []store.Relationship
	// TODO: PreLoad N + 1 ?
	if err := store.GetDB().Find(&relationships).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve relationships"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"relationships": relationships,
	})
}
