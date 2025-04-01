package controller

import (
	"fmt"
	"net/http"
	"time"

	"rental-property-management-system/backend/controller/middleware"
	"rental-property-management-system/backend/store"
	"rental-property-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// 创建订单的接口
type CreateOrderRequest struct {
	RoomID    uint `json:"room_id"`
	TotalTerm uint `json:"total_term"`
}

func CreateOrder(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var request CreateOrderRequest
	// 获取前端传来的数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 至少要租 6 个月
	if request.TotalTerm < 6 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The minimum number of terms is 6"})
		return
	}

	// 获取房间信息，检查是否已租出去
	var room store.Room
	tx := store.GetDB().Find(&room, request.RoomID)
	if tx.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": tx.Error.Error()})
		return
	}
	if tx.RowsAffected != 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room does not exist"})
		return
	}

	// 检查房间是否已被租出去
	if room.Available {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": " The room has been rented"})
		return
	}

	// 创建订单
	order := store.Order{
		UserID:           user.ID,
		RoomID:           request.RoomID,
		TotalTerm:        request.TotalTerm,
		RemainingBiilNum: request.TotalTerm,
	}
	if err := store.GetDB().Create(&order).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create the order"})
		return
	}

	// 更新房间状态为已租
	// TODO: concurrency problems
	room.Available = false
	if err := store.GetDB().Save(&room).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update room status"})
		return
	}

	// 按房间数量最少的管理员分配
	var admin store.User
	if err := store.GetDB().Where("role = ?", store.UserRoleAdmin).
		Order("managed_rooms ASC").
		First(&admin).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No available admin"})
		return
	}

	// 首次租房要付两个月押金，并在退租的时候返还
	deposit := 2 * room.Price

	// 创建支付订单
	billing := store.Billing{
		Type:   string(store.BillingTypeRentRoom),
		UserID: user.ID,
		Paid:   false,
		Price:  deposit,
		Name:   fmt.Sprintf("%s签约账单", room.Name),
	}
	if err := store.GetDB().Save(&billing).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create the bill"})
		return
	}

	// 创建 Relationship 记录
	relationship := store.Relationship{
		UserID:       user.ID,
		AdminID:      admin.ID,
		RoomID:       room.ID,
		DepositPrice: deposit,
		OrderID:      order.ID,
	}
	if err := store.GetDB().Create(&relationship).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relationship"})
		return
	}
	// 更新管理员所管房间数量
	admin.ManagedRooms++
	if err := store.GetDB().Save(&admin).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin room count"})
		return
	}

	// 为房间生成临时密码
	password := store.Password{
		RoomID:   request.RoomID,
		Password: utils.GenerateRandomPassword(6),
		// 过期时间为两个小时
		ExpiresAt: time.Now().Add(120 * time.Minute),
	}
	if err := store.GetDB().Save(&password).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password"})
		return
	}

	// 返回订单详情
	c.JSON(http.StatusOK, gin.H{
		"message": "Order created successfully",
	})
}

func ListOrders(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orders := []store.Order{}
	if err := store.GetDB().Where("id = ?", user.ID).Find(&orders).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// 支付订单接口
type PayBillRequest struct {
	BillingID uint `json:"billing_id"`
}

func PayBill(c *gin.Context) {
	var request PayBillRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找支付的订单
	billing := store.Billing{}
	if err := store.GetDB().Where("id = ?", request.BillingID).Find(&billing).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Billing not found"})
		return
	}
	if billing.Paid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Billing already paid"})
		return
	}

	// 更新订单状态
	billing.Paid = true
	if err := store.GetDB().Save(&billing).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful",
	})
}

// 用户退租接口
// TODO: 到期自动过期
type TerminateLeaseRequest struct {
	RoomID uint `json:"room_id"`
}

// TODO: 这些 SQL 应当设计成 transaction 不然真会爆炸
// TODO: 这里应该设计成创建工单，等待管理员验收后再退还押金
func TerminateLease(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var request TerminateLeaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询相关 relationship 关系表
	var relationship store.Relationship
	if err := store.GetDB().
		Where("user_id = ? AND room_id = ?", user.ID, request.RoomID).
		Find(&relationship).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Rental relationship not found"})
		return
	}

	// 更新管理员房间数量 -1
	var admin store.User
	if err := store.GetDB().Find(&admin, relationship.AdminID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	if admin.ManagedRooms > 0 {
		admin.ManagedRooms--
	}
	if err := store.GetDB().Save(&admin).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin room count"})
		return
	}

	// 提前生成当月的订单
	// WARN: 这个要和定时生成账单的协程进行一个同步
	room := store.Room{}
	if err := store.GetDB().Find(&room, request.RoomID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	billing := store.Billing{
		Type:   string(store.BillingTypeMonthlyPayment),
		Price:  utils.CalculateProRatedRent(room.Price, time.Now()),
		UserID: user.ID,
		Paid:   false,
		Name:   fmt.Sprintf("%s退租账单", room.Name),
	}
	if err := store.GetDB().Save(&billing); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create the bill"})
		return
	}

	// 更新订单状态为失效
	order := store.Order{
		ID:               relationship.OrderID,
		RemainingBiilNum: 0,
	}
	if err := store.GetDB().Save(&order); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set the order %q to invalid", relationship.OrderID)})
		return
	}

	// 生成收房工单等待管理员处理
	workOrder := store.WorkOrder{
		Type:    store.WorkOrderTypeTerminateLease,
		Status:  store.WorkOrderStatusPending,
		AdminID: relationship.AdminID,
	}
	if err := store.GetDB().Save(&workOrder); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create work order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lease terminated successfully, work order created",
	})
}

func ListBillings(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	billings := []store.Billing{}
	if err := store.GetDB().Where("user_id = ?", user.ID).Find(&billings).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"billings": billings,
	})
}
