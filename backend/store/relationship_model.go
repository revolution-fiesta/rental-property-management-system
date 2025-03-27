package store

import "gorm.io/gorm"

// Relationship 模型，表示用户与管理员和房间的关系
type Relationship struct {
	gorm.Model
	// 用户ID
	UserID uint `gorm:"not null"`
	// 管理员ID
	AdminID uint `gorm:"not null"`
	// 房间ID
	RoomID uint `gorm:"not null"`
	// 记录交付的租金
	DepositPrice float64 `gorm:"not null"`
	// 记录订单 ID
	OrderID uint `gorm:"not null"`
	// 分配时间
	AssignedAt string `gorm:"type:timestamp;"`
	// 通过外键关联房间信息
	Room Room `gorm:"foreignKey:RoomID"`
}
