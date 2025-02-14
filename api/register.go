package api

import (
	"fmt"
	"net/http"
	"rental-property-management-system/internal/database"
	"rental-property-management-system/internal/models"
	"rental-property-management-system/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

func hashPassword(password string) (string, error) {
    salt := []byte("some_random_salt") // 可以使用更复杂的盐值
    hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32) // 生成哈希值
    return fmt.Sprintf("%x", hash), nil
}

// 注册用户接口
func Register(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 判断是否已经有管理员
    var firstAdmin models.User
    if err := database.DB.Where("role = ?", "admin").First(&firstAdmin).Error; err != nil && err.Error() == "record not found" {
        // 如果没有管理员，第一位注册的用户自动成为管理员
        user.Role = "admin"
    } else {
        user.Role = "user" // 默认普通用户
    }

    // 检查用户名是否已存在
    var existingUser models.User
    if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
        return
    }

    // 密码加密
    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user.Password = hashedPassword

    // 插入用户
    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
// 管理员注册函数
func RegisterAdmin(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析请求中的 Token
	user, err := middleware.ParseToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 根据角色判断是否是管理员
	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// 密码加密
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建管理员用户
	newAdmin := models.User{
		Username:   request.Username,
		Password: hashedPassword,
		Email:      request.Email,
		Role:       "admin", // 设置角色为 admin
	}

	if err := database.DB.Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin user registered successfully"})
}