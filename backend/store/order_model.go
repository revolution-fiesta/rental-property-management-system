package store

import (
	"time"
)

type Order struct {
	// 订单ID
	ID uint `gorm:"primaryKey"`
	// 关联用户ID
	UserID uint `gorm:"not null"`
	// 房间ID
	RoomID uint `gorm:"not null"`
	// 总共期数
	TotalTerm uint `gorm:"type:decimal(10,2);not null"`
	// 创建时间
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// 更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
