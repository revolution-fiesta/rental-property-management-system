package store

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleMember UserRole = "member"
	// 用于没有注册的用户
	UserRoleNone UserRole = "none"
)

type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	Email        string `json:"email"`
	Role         string `gorm:"default:'member'"` // 区分普通用户和管理员
	ManagedRooms int    `gorm:"default:0"`        // 管理的房间数量，只有管理员有这个字段
	Salt         string `json:"salt"`
}
