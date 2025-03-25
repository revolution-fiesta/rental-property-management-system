package models

import "gorm.io/gorm"

// Relationship 模型，表示用户与管理员和房间的关系
type Relationship struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`        // 用户ID
	AdminID    uint   `gorm:"not null"`        // 管理员ID
	RoomID     uint   `gorm:"not null"`        // 房间ID
	AssignedAt string `gorm:"type:timestamp;"` // 分配时间

	// 通过外键关联房间信息
	Room Room `gorm:"foreignKey:RoomID"`
}
