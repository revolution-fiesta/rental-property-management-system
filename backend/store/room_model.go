package store

import (
	"gorm.io/gorm"
)

type RoomType string

// B: Bedroom.
// L: Living room.
const (
	RoomTypeB2L1 RoomType = "b2l1"
	RoomTypeB1L1 RoomType = "b1l1"
	RoomTypeB1   RoomType = "b1"
)

type Room struct {
	gorm.Model
	Type      RoomType `gorm:"type:varchar(20);not null"`
	Price     float64  `gorm:"type:decimal(10,2);not null"`
	Available bool     `gorm:"default:true"`
	Name      string   `gorm:"type:varchar(255)"`
	Floor     int      `gorm:"type:int;default:1;not null"`
	// 标签，最多255个字符
	Tags string `gorm:"type:varchar(255)"`
	// 占地面积，单位是平方米
	Area float64 `gorm:"type:decimal(10,2);not null"`
	// 描述信息
	Description string `gorm:"type:text"`
}
