package controller

import (
	"net/http"
	"rental-property-management-system/backend/models"
	"rental-property-management-system/backend/store"

	"time"

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
	if err := store.GetDB().First(&room, request.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在"})
		return
	}

	// 检查房间是否已被租出去
	if room.IsDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该房间已出租"})
		return
	}

	// 获取当前时间
	now := time.Now()
	year, month, _ := now.Date()
	loc := now.Location()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)      // 获取当月最后一天
	totalDays := lastOfMonth.Day()                     // 获取当月总天数
	remainingDays := lastOfMonth.Day() - now.Day() + 1 // 剩余天数（包含今天）

	// 计算本月租金
	currentMonthRent := (float64(remainingDays) / float64(totalDays)) * room.Price

	// 计算押金（2 个月）
	deposit := 2 * room.Price

	// 计算签约时需要支付的费用
	totalInitialPayment := currentMonthRent + deposit

	// 创建订单
	order := models.Order{
		UserID:     request.UserID,
		RoomID:     request.RoomID,
		Status:     models.Pending, // 订单状态：待支付
		TotalPrice: totalInitialPayment,
		//StartDate:  now, // 记录租赁开始时间
	}

	// 保存订单到数据库
	if err := store.GetDB().Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 更新房间状态为已租
	room.IsDeleted = true
	store.GetDB().Save(&room)

	// 返回订单详情
	c.JSON(http.StatusOK, gin.H{
		"message":       "订单创建成功，请支付租金和押金",
		"order_id":      order.ID,
		"current_rent":  currentMonthRent,
		"deposit":       deposit,
		"total_payment": totalInitialPayment,
	})
}

// 生成每月月结订单
// 生成月结订单
func GenerateMonthlyOrders(c *gin.Context) {
	now := time.Now()
	// 查询所有正在租赁的订单
	var activeOrders []models.Order
	if err := store.GetDB().Where("status = ?", models.Completed).Find(&activeOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单失败"})
		return
	}

	var newOrders []models.Order

	for _, order := range activeOrders {
		var room models.Room
		if err := store.GetDB().First(&room, order.RoomID).Error; err != nil {
			continue
		}

		// 计算当月租金
		monthlyRent := room.Price

		// 生成新的月结订单
		newOrder := models.Order{
			UserID:     order.UserID,
			RoomID:     order.RoomID,
			Status:     models.Pending,
			TotalPrice: monthlyRent,
			CreatedAt:  now,
		}
		newOrders = append(newOrders, newOrder)
	}

	// 批量插入订单
	if len(newOrders) > 0 {
		if err := store.GetDB().Create(&newOrders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建订单失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "月结订单已生成", "order_count": len(newOrders)})
}

// 支付订单接口
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
	if err := store.GetDB().First(&order, request.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == models.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already paid"})
		return
	}

	// 更新订单状态
	order.Status = models.Completed
	if err := store.GetDB().Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// 获取房间信息，并更新状态为已租出
	var room models.Room
	if err := store.GetDB().First(&room, order.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	room.IsDeleted = true // 表示已租
	if err := store.GetDB().Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room status"})
		return
	}

	// 按房间数量最少的管理员分配
	var admin models.User
	if err := store.GetDB().Where("role = ?", "admin").
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

	if err := store.GetDB().Create(&relationship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relationship"})
		return
	}

	// 更新管理员所管房间数量
	admin.ManagedRooms++
	if err := store.GetDB().Save(&admin).Error; err != nil {
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
	if err := store.GetDB().
		Where("user_id = ? AND room_id = ?", request.UserID, request.RoomID).
		First(&relationship).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rental relationship not found"})
		return
	}

	// 查询房间信息
	var room models.Room
	if err := store.GetDB().First(&room, request.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// 查询管理员信息
	var admin models.User
	if err := store.GetDB().First(&admin, relationship.AdminID).Error; err != nil {
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
	if err := store.GetDB().Delete(&relationship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relationship"})
		return
	}

	// 更新房间状态为未租出
	room.IsDeleted = false
	if err := store.GetDB().Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room status"})
		return
	}

	// 更新管理员房间数量 -1
	if admin.ManagedRooms > 0 {
		admin.ManagedRooms--
	}
	if err := store.GetDB().Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin room count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Rental cancelled successfully, deposit refunded",
		"refund_amount": refundAmount,
	})
}
