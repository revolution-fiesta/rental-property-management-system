package controllers

import (
	"net/http"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"

	"github.com/gin-gonic/gin"
)

// 创建订单的接口
func CreateOrder(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id"`
		RoomID uint `json:"room_id"`
	}

	// 获取前端传来的数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取房间信息，检查是否已租出去
	var room models.Room
	if err := database.DB.First(&room, request.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 检查房间是否已被租出去
	if room.IsDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room is already rented"})
		return
	}

	// 假设前端会传入起租月数，至少6个月
	var rentRequest struct {
		Months int `json:"months"`
	}
	if err := c.ShouldBindJSON(&rentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent duration"})
		return
	}
	if rentRequest.Months < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum rental period is 6 months"})
		return
	}

	// 计算订单总价
	totalPrice := float64(rentRequest.Months + 2) * room.Price

	// 创建订单
	order := models.Order{
		UserID:    request.UserID,
		RoomID:    request.RoomID,
		Status:    models.Pending,
		TotalPrice: totalPrice,
	}

	// 保存订单到数据库
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// 模拟调用支付窗口的逻辑（这里前端会调用微信支付接口）
	c.JSON(http.StatusOK, gin.H{
		"message":    "Order created successfully, please proceed to payment",
		"order_id":   order.ID,
		"total_price": totalPrice,
	})
}
func PayOrder(c *gin.Context) {
	var request struct {
		OrderID uint `json:"order_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找订单
	var order models.Order
	if err := database.DB.First(&order, request.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == models.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already paid"})
		return
	}

	// 更新订单状态
	order.Status = models.Completed
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// 获取房间信息，并更新状态为已租出
	var room models.Room
	if err := database.DB.First(&room, order.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	room.IsDeleted = true // 表示已租
	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room status"})
		return
	}

	// 按房间数量最少的管理员分配
	var admin models.User
	if err := database.DB.Where("role = ?", "admin").
		Order("room_count ASC").
		First(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No available admin"})
		return
	}

	// 创建 Relationship 记录
	relationship := models.Relationship{
		UserID:  order.UserID,
		AdminID: admin.ID,
		RoomID:  order.RoomID,
	}

	if err := database.DB.Create(&relationship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relationship"})
		return
	}

	// 更新管理员所管房间数量
	admin.ManagedRooms++
	if err := database.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin room count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful, room rented, admin assigned",
		"order":   order,
	})
}
// 用户退租接口
func CancelRental(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id"`
		RoomID uint `json:"room_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询 relationship 关系表，找到对应的记录
	var relationship models.Relationship
	if err := database.DB.
		Where("user_id = ? AND room_id = ?", request.UserID, request.RoomID).
		First(&relationship).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rental relationship not found"})
		return
	}

	// 查询房间信息
	var room models.Room
	if err := database.DB.First(&room, request.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 查询管理员信息
	var admin models.User
	if err := database.DB.First(&admin, relationship.AdminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// 退还押金（这里可以根据业务逻辑调整，暂时打印模拟退款）
	// 实际项目应调用支付平台接口进行退款，这里仅作业务流程展示
	// 假设押金是2个月房租
	refundAmount := room.Price * 2
	// 模拟退款操作
	c.JSON(http.StatusOK, gin.H{"message": "Deposit refund initiated", "refund_amount": refundAmount})

	// 删除 relationship 表中的记录
	if err := database.DB.Delete(&relationship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relationship"})
		return
	}

	// 更新房间状态为未租出
	room.IsDeleted = false
	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room status"})
		return
	}

	// 更新管理员房间数量 -1
	if admin.ManagedRooms > 0 {
		admin.ManagedRooms--
	}
	if err := database.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin room count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rental cancelled successfully, deposit refunded",
		"refund_amount": refundAmount,
	})
}