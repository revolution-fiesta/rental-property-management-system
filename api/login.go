package api

import (
	"fmt"
	"net/http"

	"rental-property-management-system/middleware"
	"rental-property-management-system/models"
	"rental-property-management-system/store"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

func verifyPassword(storedHash, password string) bool {
	salt := []byte("some_random_salt") // 必须使用相同的盐值
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hash) == storedHash
}

// 生成 token 的例子
func GenerateToken(username, role string) (string, error) {
	// 创建 token
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // token 过期时间（设置为 72 小时后过期）
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥生成 token
	jwtSecret := []byte(middleware.GenerateRandomKey())
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 登录接口
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := store.GetDB().Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	// 验证密码
	if !verifyPassword(user.Password, loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// 生成JWT token
	token, err := GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"role":    user.Role, // 返回用户角色，后续可以根据角色做权限验证
	})
}
