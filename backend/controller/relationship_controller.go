package controller

import (
	"net/http"

	"rental-property-management-system/backend/controller/middleware"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

// 根据用户 ID 查询房间分配关系
// TODO: PreLoad N + 1 ?
// TODO: 限制权限
func ListRelationships(c *gin.Context) {
	var relationships []store.Relationship
	if err := store.GetDB().Find(&relationships).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve relationships"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"relationships": relationships,
	})
}

func ListOwnedRooms(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	relationships := []store.Relationship{}
	tx := store.GetDB().Where("user_id = ?", user.ID).Find(&relationships)
	if err := tx.Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rooms := []store.Room{}
	for _, relationship := range relationships {
		room := store.Room{}
		tx := store.GetDB().Where("id = ?", relationship.RoomID).Find(&room)
		if err := tx.Error; err != nil || tx.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}
