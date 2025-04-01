package store

import "gorm.io/gorm"

type BillingType string

const (
	BillingTypeRentRoom       BillingType = "rent-room"
	BillingTypeMonthlyPayment BillingType = "monthly-pay"
	BillingTypeTerminateLease BillingType = "terminate-lease"
)

// Relationship 模型，表示用户与管理员和房间的关系
type Billing struct {
	// WARN: 这里的 ID 和 billingID 有点问题
	gorm.Model
	Name   string  `gorm:"not null"`
	Type   string  `gorm:"not null"`
	UserID uint    `gorm:"not null"`
	Price  float64 `gorm:"not null"`
	Paid   bool    `gorm:"default:false"`
}
