package store

import (
	"gorm.io/gorm"
)

type Signature struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// 用于存储签名图片的 Base64 数据
	ImageData string `json:"image_data"`
}
