package controller

import (
	"net/http"
	"regexp"
	"rental-property-management-system/backend/store"
	"time"

	"github.com/gin-gonic/gin"
)

// 修改房间密码
type ChangeRoomPasswordRequest struct {
	RoomID      uint   `json:"room_id" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ChangeRoomPassword(c *gin.Context) {
	var request ChangeRoomPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 验证新密码格式（必须是6位数字）
	if !isValidPassword(request.NewPassword) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "New password must be a 6-digit number"})
		return
	}

	// 查询需要更新密码的房间
	password := store.Password{}
	if err := store.GetDB().First(&password, "room_id = ?", request.RoomID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 设置有效期为永久并更新密码
	password.Password = request.NewPassword
	password.ExpiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	if err := store.GetDB().Save(&password).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// 根据 ID 获取房间密码
type GetPasswordRequest struct {
	RoomID string `json:"room_id"`
}

func GetPassword(c *gin.Context) {
	var request GetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	password := store.Password{}
	if err := store.GetDB().Where("room_id = ?", request.RoomID).Find(&password).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password": password,
	})
}

// 校验密码格式（6位数字）
func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^\d{6}$`)
	return re.MatchString(password)
}
