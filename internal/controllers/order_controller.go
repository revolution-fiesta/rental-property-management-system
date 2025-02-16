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

// 创建订单的API
func CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查房间是否存在
	var room models.Room
	if err := database.DB.Where("id = ?", request.RoomID).First(&room).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 创建订单
	order := models.Order{
		UserID:    request.UserID,
		RoomID:    request.RoomID,
		Status:    models.Pending, // 默认订单状态是待处理
		TotalPrice: request.TotalPrice,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// 这里调用支付平台（如微信支付、支付宝等），返回支付信息
	paymentInfo := generatePaymentInfo(order) // 生成支付信息

	// 返回支付信息
	c.JSON(http.StatusOK, gin.H{
		"message": "Order created successfully",
		"order":   order,
		"payment_info": paymentInfo, // 返回支付信息
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
// 删除订单的API
func DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	// 查找并删除订单
	if err := database.DB.Delete(&models.Order{}, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
// 支付回调的API
func PaymentCallback(c *gin.Context) {
	orderID := c.DefaultQuery("order_id", "")  // 支付平台传递的订单ID
	status := c.DefaultQuery("status", "")      // 支付状态（例如 "success" 或 "fail"）

	// 查找订单
	var order models.Order
	if err := database.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 根据支付状态更新订单
	if status == "success" {
		order.Status = models.Confirmed // 支付成功，更新状态为已确认
	} else {
		order.Status = models.Cancelled // 支付失败，更新状态为已取消
	}

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}