package store

import (
	"gorm.io/gorm"
)

type RoomType string

const (
	TwoBedroom RoomType = "two_bedroom"
	OneBedroom RoomType = "one_bedroom"
	SingleRoom RoomType = "single_room"
)

type Room struct {
	gorm.Model
	Type      RoomType `gorm:"type:varchar(20);not null"`
	Quantity  int      `gorm:"not null"`
	Price     float64  `gorm:"type:decimal(10,2);not null"`
	IsDeleted bool     `gorm:"default:false"`
	Tags      string   `gorm:"type:varchar(255)"`           // 标签，最多255个字符
	Area      float64  `gorm:"type:decimal(10,2);not null"` // 占地面积，单位是平方米
}
