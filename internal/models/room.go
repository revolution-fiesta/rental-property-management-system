package models

import (
	"gorm.io/gorm"
)

type RoomType string

const (
	TwoBedroom   RoomType = "two_bedroom"
	OneBedroom   RoomType = "one_bedroom"
	SingleRoom   RoomType = "single_room"
)

type Room struct {
	gorm.Model
	Type      RoomType `gorm:"type:varchar(20);not null"`
	Quantity  int      `gorm:"not null"`
	Price     float64  `gorm:"type:decimal(10,2);not null"`
	IsDeleted bool     `gorm:"default:false"`
}