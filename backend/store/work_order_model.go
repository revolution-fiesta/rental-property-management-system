package store

import (
	"time"
)

type WorkOrderStatus string

const (
	WorkOrderPending   WorkOrderStatus = "pending"    // 待处理
	WorkOrderInProcess WorkOrderStatus = "in_process" // 处理中
	WorkOrderCompleted WorkOrderStatus = "completed"  // 已完成
)

type WorkOrder struct {
	ID        uint            `gorm:"primaryKey"`
	UserID    uint            `gorm:"not null"` // 报修的用户ID
	User      User            `gorm:"foreignKey:UserID"`
	RoomID    uint            `gorm:"not null"` // 报修的房间ID
	Room      Room            `gorm:"foreignKey:RoomID"`
	AdminID   uint            `gorm:"not null"` // 处理的管理员ID
	Admin     User            `gorm:"foreignKey:AdminID"`
	Problem   string          `gorm:"type:text;not null"` // 问题描述，例如：漏水、门锁坏了
	Status    WorkOrderStatus `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
}
