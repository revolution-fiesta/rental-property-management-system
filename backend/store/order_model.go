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
	// 剩余需要生成账单的期数, 如果为 0 则说明订单已失效
	RemainingBiilNum uint `gorm:"not null"`
	// 创建时间
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// 更新时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
