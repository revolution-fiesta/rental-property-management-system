package store

import (
	"time"
)

type WorkOrderStatus string

const (
	// 待处理
	WorkOrderPending WorkOrderStatus = "pending"
	// 处理中
	WorkOrderInProcess WorkOrderStatus = "in_process"
	// 已完成
	WorkOrderCompleted WorkOrderStatus = "completed"
)

type WorkOrderProblemType string

const (
	// TODO: 能不能不用中文
	WorkOrderProblemTerminateLease = "退租验收"
)

type WorkOrder struct {
	ID uint `gorm:"primaryKey"`
	// 报修的用户ID
	UserID uint `gorm:"not null"`
	// 报修的房间ID
	RoomID uint `gorm:"not null"`
	// 问题描述，例如：漏水、门锁坏了
	Problem   string          `gorm:"type:text;not null"`
	Status    WorkOrderStatus `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
}
