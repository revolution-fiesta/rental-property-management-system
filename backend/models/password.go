package models

import (
	"time"
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	RoomID    uint      `gorm:"not null"`
	Password  string    `gorm:"type:varchar(20);not null"`
	IsTemp    bool      `gorm:"default:true"`
	ExpiresAt time.Time `gorm:"not null"`
}