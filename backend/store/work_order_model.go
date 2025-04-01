package store

import (
	"time"
)

type WorkOrderStatus string

const (
	// 待处理
	WorkOrderStatusPending WorkOrderStatus = "pending"
	// 处理中
	WorkOrderStatusInProcess WorkOrderStatus = "in_process"
	// 已完成
	WorkOrderStatusCompleted WorkOrderStatus = "completed"
)

type WorkOrderType string

const (
	// TODO: 能不能不用中文
	WorkOrderTypeTerminateLease = "退租验收"
	WorkOrderTypeGeneral        = "一般维护"
)

type WorkOrder struct {
	ID uint `gorm:"primaryKey"`
	// 报修的用户ID
	UserID uint `gorm:"not null"`
	// 报修的房间ID
	RoomID uint `gorm:"not null"`
	// 工单类型
	Type string `gorm:"type:text;not null"`
	// 问题描述
	Description string `gorm:"type:text"`
	// 关联的管理员 ID
	AdminID   uint            `gorm:"not null"`
	Status    WorkOrderStatus `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
}
