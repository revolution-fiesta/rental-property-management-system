package controllers

import (
	"net/http"
	"fmt"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"
	"github.com/gin-gonic/gin"
)

// 创建订单请求数据结构
type CreateOrderRequest struct {
	UserID    uint    `json:"user_id" binding:"required"`
	RoomID    uint    `json:"room_id" binding:"required"`
	TotalPrice float64 `json:"total_price" binding:"required"`
}
// 更新订单状态请求数据结构
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
// 模拟生成支付信息函数（实际使用时需要集成支付平台SDK）
func generatePaymentInfo(order models.Order) string {
	// 这里你需要根据支付平台的API生成支付链接或支付请求参数
	// 以下为伪代码示例

	// 假设支付平台需要传递的参数
	// 这里只是一个简单示范，实际情况需要根据支付平台的文档进行实现
	paymentLink := fmt.Sprintf("https://paymentgateway.com/pay?order_id=%d&amount=%.2f", order.ID, order.TotalPrice)
	return paymentLink
}

// 创建订单的 API
func CreateOrder(c *gin.Context) {
	var request struct {
		UserID     uint  `json:"user_id"`
		RoomID     uint  `json:"room_id"`
		RentMonths int   `json:"rent_months"` // 用户填写的租期
	}

	// 获取前端传来的数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证租期是否至少6个月
	if request.RentMonths < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rent months must be at least 6"})
		return
	}

	// 获取房间信息
	var room models.Room
	if err := database.DB.First(&room, request.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 计算总租金
	totalPrice := float64(request.RentMonths + 2) * room.Price

	// 创建订单
	order := models.Order{
		UserID:    request.UserID,
		RoomID:    request.RoomID,
		Status:    models.Pending,  // 默认订单状态为待处理
		TotalPrice: totalPrice,
	}

	// 保存订单到数据库
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

// 更新订单状态的API
func UpdateOrderStatus(c *gin.Context) {
	var request UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取订单ID
	orderID := c.Param("id")

	// 查找订单
	var order models.Order
	if err := database.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 更新订单状态
	order.Status = models.OrderStatus(request.Status)
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "order": order})
}
// 支付完订单后更新房间状态
func PayOrder(c *gin.Context) {
	var request struct {
		OrderID uint `json:"order_id"`
	}

	// 获取前端传来的数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取订单信息
	var order models.Order
	if err := database.DB.First(&order, request.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 更新订单为已支付
	order.Status = models.Completed
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// 获取房间信息并更新状态
	var room models.Room
	if err := database.DB.First(&room, order.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 将房间状态更新为已租
	room.IsDeleted = true
	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room status"})
		return
	}

	// 查找房间数量最少的管理员
	var admin models.User
	if err := database.DB.Where("role = ?", "admin").
		Order("room_count ASC").First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No admin found"})
		return
	}

	// 创建管理员与用户、房间的关系
	relationship := models.Relationship{
		UserID:    order.UserID,
		AdminID:   admin.ID,
		RoomID:    order.RoomID,
		AssignedAt: fmt.Sprintf("%s", order.CreatedAt),
	}

	// 保存关系数据
	if err := database.DB.Create(&relationship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign room to user"})
		return
	}

	// 更新管理员的房间数量
	admin.ManagedRooms++
	if err := database.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin's room count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful, room rented, and admin assigned",
		"order":   order,
	})
}