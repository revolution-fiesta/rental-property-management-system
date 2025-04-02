package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"rental-property-management-system/backend/store"
	"rental-property-management-system/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 修改房间密码
type ChangeRoomPasswordRequest struct {
	RoomID      uint   `json:"room_id"`
	NewPassword string `json:"new_password"`
}

// WARN: 需要验证用户权限
// TODO: 能不能用一条 gorm 语句重构
func ChangeRoomPassword(c *gin.Context) {
	var request ChangeRoomPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	password := store.Password{}
	// 查询需要更新密码的房间
	if err := store.GetDB().First(&password, "room_id = ?", request.RoomID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 如果没有传递密码参数则随机生成新密码,
	// 否则验证新密码格式（必须是6位数字）
	if request.NewPassword == "" {
		password.Password = utils.GenerateRandomPassword(6)
		password.ExpiresAt = time.Now().Add(120 * time.Minute)
	} else {
		if !isValidPassword(request.NewPassword) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "New password must be a 6-digit number"})
			return
		}
		// 设置有效期为永久并更新密码
		password.Password = request.NewPassword
		password.ExpiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}

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
	RoomID uint `json:"room_id"`
}

func GetPassword(c *gin.Context) {
	var request GetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	password := store.Password{}
	tx := store.GetDB().Where("room_id = ?", request.RoomID).Find(&password)
	if err := tx.Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if tx.RowsAffected != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to get password from room id: %d", request.RoomID),
		})
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
