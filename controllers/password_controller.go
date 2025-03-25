package controllers

import (
	"math/rand"
	"net/http"
	"regexp"
	"rental-property-management-system/models"
	"rental-property-management-system/store"

	"time"

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
		ExpiresAt: time.Now().Add(120 * time.Minute),
	}

	if err := store.GetDB().Create(&password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"password":   password.Password,
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

// 修改房间密码
func changeRoomPassword(c *gin.Context) {
	// 请求结构体定义
	var request struct {
		RoomID      uint   `json:"room_id" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// 解析请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 验证新密码格式（必须是6位数字）
	if !isValidPassword(request.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be a 6-digit number"})
		return
	}

	// 查询该房间的临时密码记录
	var password models.Password
	if err := store.GetDB().First(&password, "room_id = ? AND is_temp = ?", request.RoomID, true).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found or no temporary password available"})
		return
	}

	// 验证旧密码
	if password.Password != request.OldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	// 更新密码并设置有效期为永久
	password.Password = request.NewPassword
	password.ExpiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC) // 设置密码有效期为永久

	// 保存更新后的密码
	if err := store.GetDB().Save(&password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Password updated successfully",
		"new_password": password.Password,
		"expires_at":   password.ExpiresAt,
	})
}

// 校验密码格式（6位数字）
func isValidPassword(password string) bool {
	// 使用正则表达式验证密码是否为6位数字
	re := regexp.MustCompile(`^\d{6}$`)
	return re.MatchString(password)
}
