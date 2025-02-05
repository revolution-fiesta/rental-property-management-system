package controllers

import (
	"net/http"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"
	
	"github.com/gin-gonic/gin"
	
)
// 初始化房间数据
func InitRoomData() {
	rooms := []models.Room{
		{Type: models.TwoBedroom, Quantity: 311, Price: 5000},   // 两房一厅
		{Type: models.OneBedroom, Quantity: 605, Price: 3500},   // 一房一厅
		{Type: models.SingleRoom, Quantity: 505, Price: 2000},   // 单间
	}

	for _, room := range rooms {
		database.DB.FirstOrCreate(&room, models.Room{Type: room.Type})
	}
}
func GetAvailableRooms(c *gin.Context) {
	var rooms []models.Room

	result := database.DB.Where("is_deleted = ?", false).Find(&rooms)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rooms})
}
// 选择房间
func selectRoom(c *gin.Context) {
	type RoomSelection struct {
		RoomType models.RoomType `json:"room_type" binding:"required"`
		Quantity int      `json:"quantity" binding:"required"`
	}

	var selection RoomSelection

	// 解析用户请求的房间类型和数量
	if err := c.ShouldBindJSON(&selection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	var room models.Room

	// 根据房间类型和库存筛选房间
	if err := database.DB.First(&room, "type = ? AND is_deleted = ? AND quantity >= ?", selection.RoomType, false, selection.Quantity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间类型不可用或库存不足"})
		return
	}

	// 更新房间数量
	room.Quantity -= selection.Quantity
	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新房间库存"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "房间选择成功",
		"room":    room,
	})
}

