package store

import (
	"time"

	"gorm.io/gorm"
)

// TODO: 也许什么时候设置一下边缘设备
// 数据库中的信息会通过某种方式下放到边缘设备
// 更新密码后，新密码和过期时间都会同步到边缘设备中
type Password struct {
	gorm.Model
	RoomID    uint      `gorm:"not null"`
	Password  string    `gorm:"type:varchar(20);not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
