package models


type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Password string `json:"password"`
	Email        string `json:"email"`
	Role     string `gorm:"default:'user'"` // 区分普通用户和管理员
}