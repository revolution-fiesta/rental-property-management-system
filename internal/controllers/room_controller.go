package controllers

import (
	"net/http"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"fmt"
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
// 获取可用的房间
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
// 管理员更新房间的信息
func UpdateRoomInfo(c *gin.Context) {
	// 通过中间件获取管理员权限
	user, _ := c.Get("user") // 获取用户信息

	// 确认用户为管理员
	if user == nil || user.(*models.User).Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
		return
	}

	// 定义请求体结构
	var request struct {
		RoomID    uint     `json:"room_id" binding:"required"`
		Type      *string  `json:"type"`
		Price     *float64 `json:"price"`
		IsDeleted *bool    `json:"is_deleted"`
		Tags      *string  `json:"tags"`
		Area      *float64 `json:"area"`
	}

	// 绑定请求体数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 查找房间
	var room models.Room
	err := database.DB.First(&room, request.RoomID).Error

	// 如果找不到房间，则创建新的房间
	if err != nil && err == gorm.ErrRecordNotFound {
		// 房间不存在，创建新房间
		room = models.Room{
			Type:      models.RoomType(*request.Type),
			Price:     *request.Price,
			IsDeleted: *request.IsDeleted,
			Tags:      *request.Tags,
			Area:      *request.Area,
		}
		if err := database.DB.Create(&room).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Room created successfully",
			"room":    room,
		})
		return
	}

	// 如果房间存在，则更新房间信息
	if request.Type != nil {
		room.Type = models.RoomType(*request.Type)
	}
	if request.Price != nil {
		room.Price = *request.Price
	}
	if request.IsDeleted != nil {
		room.IsDeleted = *request.IsDeleted
	}
	if request.Tags != nil {
		room.Tags = *request.Tags
	}
	if request.Area != nil {
		room.Area = *request.Area
	}

	// 更新数据库
	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room updated successfully",
		"room":    room,
	})
}
// 查询所有房间接口
func GetAllRooms(c *gin.Context) {
	var rooms []models.Room

	// 查询所有房间数据
	if err := database.DB.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	// 返回房间信息
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}
// 查询房间信息，支持多个过滤条件
func GetFilteredRooms(c *gin.Context) {
	// 获取查询参数
	priceMin := c.DefaultQuery("price_min", "0")
	priceMax := c.DefaultQuery("price_max", "1000000000")
	roomType := c.DefaultQuery("type", "")
	orientation := c.DefaultQuery("orientation", "")
	areaMin := c.DefaultQuery("area_min", "0")  // 占地面积最小值
	areaMax := c.DefaultQuery("area_max", "1000000000") // 占地面积最大值

	var rooms []models.Room
	query := database.DB

	// 按房价范围过滤
	query = query.Where("price >= ? AND price <= ?", priceMin, priceMax)

	// 按房间类型过滤
	if roomType != "" {
		query = query.Where("type = ?", roomType)
	}

	// 按朝向过滤
	if orientation != "" {
		query = query.Where("tags LIKE ?", fmt.Sprintf("%%%s%%", orientation)) // 朝向在 tags 字段中
	}

	// 按占地面积过滤
	query = query.Where("area >= ? AND area <= ?", areaMin, areaMax)

	// 查询符合条件的房间数据
	if err := query.Find(&rooms).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

