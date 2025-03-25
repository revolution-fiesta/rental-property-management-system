package models

import (
	"time"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"   // 待处理
	Confirmed OrderStatus = "confirmed" // 已确认
	Cancelled OrderStatus = "cancelled" // 已取消
	Completed OrderStatus = "completed" // 已完成
)

type Order struct {
	ID        uint        `gorm:"primaryKey"`               // 订单ID
	UserID    uint        `gorm:"not null"`                // 用户ID
	User      User        `gorm:"foreignKey:UserID"`       // 关联用户
	RoomID    uint        `gorm:"not null"`                // 房间ID
	Room      Room        `gorm:"foreignKey:RoomID"`       // 关联房间
	Status    OrderStatus `gorm:"type:varchar(20);not null"` // 订单状态
	TotalPrice float64    `gorm:"type:decimal(10,2);not null"` // 订单总价
	CreatedAt time.Time   `gorm:"autoCreateTime"`          // 创建时间
	UpdatedAt time.Time   `gorm:"autoUpdateTime"`          // 更新时间
}