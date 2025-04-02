package controller

import (
	"net/http"
	"rental-property-management-system/backend/controller/middleware"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

// 创建工单
type CreateWorkOrderRequest struct {
	RoomID      uint   `json:"room_id"`
	Description string `json:"description"`
}

func CreateWorkOrder(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var request CreateWorkOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 确认房间和用户的绑定关系，找到对应的管理员
	var relationship store.Relationship
	tx := store.GetDB().Where("user_id = ? AND room_id = ?", user.ID, request.RoomID).Find(&relationship)
	if err := tx.Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tx.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to find any related rooms"})
		return
	}

	// 创建工单
	workOrder := store.WorkOrder{
		UserID:      user.ID,
		RoomID:      request.RoomID,
		Type:        store.WorkOrderTypeGeneral,
		Description: request.Description,
		Status:      store.WorkOrderStatusPending,
		AdminID:     relationship.AdminID,
	}
	if err := store.GetDB().Create(&workOrder).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "工单创建成功",
		"work_order": workOrder,
	})
}

// 用户查询提交的工单
func ListUserWorkOrders(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}
	var workOrders []store.WorkOrder
	if err := store.GetDB().Where("user_id = ?", user.ID).Find(&workOrders).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get work orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"work_orders": workOrders})
}

// 管理员查询待处理工单
func ListAdminWorkOrders(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}
	var workOrders []store.WorkOrder
	if err := store.GetDB().Where("admin_id = ?", user.ID).Find(&workOrders).Error; err != nil {
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
	// TODO: 退租之后还涉及电子门锁的管理, 后续需要考虑
	if workOrder.Type == store.WorkOrderTypeTerminateLease && request.Status == string(store.WorkOrderStatusCompleted) {
		// 根据房间 ID 寻找绑定关系
		relationship := store.Relationship{}
		// TODO: 细致地错误处理
		if tx := store.GetDB().Where("room_id = ?", workOrder.RoomID).Find(&relationship); tx.Error != nil || tx.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "relationship not found"})
			return
		}
		// WARN: 这个是软删除，不是一次性全删完
		if err := store.GetDB().Where("id = ?", relationship.ID).Delete(&relationship).Error; err != nil {
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
