package store

import (
	"crypto/rand"
	"fmt"
	"rental-property-management-system/backend/utils"

	"github.com/pkg/errors"
)

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
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
	// 区分普通用户和管理员
	Role string `gorm:"default:'member'"`
	// 管理的房间数量，只有管理员有这个字段
	ManagedRooms int    `gorm:"default:0"`
	Salt         string `json:"salt"`
	// Used for wechat login.
	OpenID string `json:"open_id"`
}

func CreateUser(username, password, email, role, openID string) error {
	// 生成盐值与密码
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	saltString := fmt.Sprintf("%x", salt)
	hashedPasswd := utils.Sha256(password, saltString)
	// 插入用户
	newUser := User{
		Username:     username,
		PasswordHash: hashedPasswd,
		Email:        email,
		Role:         role,
		Salt:         saltString,
		OpenID:       openID,
	}
	if err := GetDB().Create(&newUser).Error; err != nil {
		return errors.Wrapf(err, "failed to register user")
	}
	return nil
}
