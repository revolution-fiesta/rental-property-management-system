package controller

import (
	"net/http"
	"rental-property-management-system/backend/controller/middleware"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 确认房间和用户的绑定关系，找到对应的管理员
	var relationship store.Relationship
	if err := store.GetDB().Where("user_id = ? AND room_id = ?", input.UserID, input.RoomID).First(&relationship).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "未找到租赁关系"})
		return
	}

	// 创建工单
	workOrder := store.WorkOrder{
		UserID: input.UserID,
		RoomID: input.RoomID,
		Type:   input.Problem,
		Status: store.WorkOrderStatusPending,
	}

	if err := store.GetDB().Create(&workOrder).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "工单创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单创建成功", "work_order": workOrder})
}

// 管理员查询待处理工单
func ListWorkOrders(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}
	var workOrders []store.WorkOrder
	if err := store.GetDB().Where("admin_id = ? AND status = ?", user.ID, store.WorkOrderStatusPending).Find(&workOrders).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get work orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"work_orders": workOrders})
}

// 管理员处理完维修后，点击“完成工单”
type UpdateWorkOrderRequest struct {
	WorkOrderID uint   `json:"work_order_id"`
	Status      string `json:"status"`
}

func UpdateWorkOrder(c *gin.Context) {
	// 获取参数并检查
	var request UpdateWorkOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid arguments"})
		return
	}
	if request.Status != string(store.WorkOrderStatusInProcess) && request.Status != string(store.WorkOrderStatusCompleted) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "非法状态"})
		return
	}

	// 获取目标工单
	var workOrder store.WorkOrder
	if err := store.GetDB().Where("id = ?", request.WorkOrderID).Find(&workOrder).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "work order does not exist"})
		return
	}

	// 根据工单类型触发对应流程
	// 退租流程
	if workOrder.Type == store.WorkOrderTypeTerminateLease && request.Status == string(store.WorkOrderStatusCompleted) {
		// 根据房间 ID 寻找绑定关系
		relationship := store.Relationship{}
		// WARN: 如果查找不到会返回错误吗
		if err := store.GetDB().Where("room_id = ?", workOrder.RoomID).Find(&relationship).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "relationship not found"})
			return
		}
		if err := store.GetDB().Delete(&relationship).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relationship"})
			return
		}
		// 更新房间状态为未租出
		var room store.Room
		if err := store.GetDB().Find(&room, workOrder.RoomID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		room.Available = true
		if err := store.GetDB().Save(&room).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room status"})
			return
		}
	}

	// 更新工单状态
	workOrder.Status = store.WorkOrderStatus(request.Status)
	if err := store.GetDB().Save(&workOrder).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to save work order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "work-order updated successfully"})
}
