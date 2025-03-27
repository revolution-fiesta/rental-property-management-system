package controller

import (
	"net/http"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

// 创建工单
func CreateWorkOrder(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"user_id"`
		RoomID  uint   `json:"room_id"`
		Problem string `json:"problem"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 确认房间和用户的绑定关系，找到对应的管理员
	var relationship store.Relationship
	if err := store.GetDB().Where("user_id = ? AND room_id = ?", input.UserID, input.RoomID).First(&relationship).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到租赁关系"})
		return
	}

	// 创建工单
	workOrder := store.WorkOrder{
		UserID:  input.UserID,
		RoomID:  input.RoomID,
		Problem: input.Problem,
		Status:  store.WorkOrderPending,
	}

	if err := store.GetDB().Create(&workOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "工单创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单创建成功", "work_order": workOrder})
}

// 管理员查询待处理工单
func ListWorkOrders(c *gin.Context) {
	// 通过中间件获取管理员权限
	user, _ := c.Get("user") // 获取用户信息

	// 确认用户为管理员
	if user == nil || user.(*store.User).Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
		return
	}
	adminID := c.Param("admin_id")
	var workOrders []store.WorkOrder
	if err := store.GetDB().Where("admin_id = ? AND status = ?", adminID, store.WorkOrderPending).Find(&workOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询工单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"work_orders": workOrders})
}

// 管理员处理完维修后，点击“完成工单”
func UpdateWorkOrderStatus(c *gin.Context) {
	// 通过中间件获取管理员权限
	user, _ := c.Get("user") // 获取用户信息

	// 确认用户为管理员
	if user == nil || user.(*store.User).Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
		return
	}
	var input struct {
		WorkOrderID uint   `json:"work_order_id"`
		Status      string `json:"status"` // "in_process", "completed"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var workOrder store.WorkOrder
	if err := store.GetDB().First(&workOrder, input.WorkOrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	if input.Status != string(store.WorkOrderInProcess) && input.Status != string(store.WorkOrderCompleted) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法状态"})
		return
	}

	workOrder.Status = store.WorkOrderStatus(input.Status)
	store.GetDB().Save(&workOrder)

	c.JSON(http.StatusOK, gin.H{"message": "工单状态已更新"})
}
