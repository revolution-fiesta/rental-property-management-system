package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	WechatOpenID string    `gorm:"type:varchar(100);uniqueIndex"`
	Phone        string    `gorm:"type:varchar(20);not null"`
	SignedUpAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}