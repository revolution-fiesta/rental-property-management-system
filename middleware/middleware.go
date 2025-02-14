package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"rental-property-management-system/internal/models"
	"github.com/gin-gonic/gin"
	 "crypto/rand"
    "encoding/base64"
    "log"
)
func GenerateRandomKey() string {
    // 生成 32 字节的随机密钥
    secret := make([]byte, 32)
    _, err := rand.Read(secret)
    if err != nil {
        log.Fatal("Failed to generate random key:", err)
    }
    // 将随机密钥编码为 Base64 字符串
    return base64.StdEncoding.EncodeToString(secret)
}
var jwtSecret = []byte(GenerateRandomKey()) // JWT 秘钥

// 解析 Token 并提取角色信息
func ParseToken(c *gin.Context) (*models.User, error) {
	// 从请求头获取 token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, fmt.Errorf("Authorization token missing or invalid")
	}

	tokenString = tokenString[7:] // 去掉 "Bearer " 前缀

	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保使用的是 HMAC 签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token: %v", err)
	}

	// 获取 Claims 中的用户信息
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		role := claims["role"].(string)

		// 返回用户信息
		return &models.User{
			Username: username,
			Role:     role,
		}, nil
	}

	return nil, fmt.Errorf("Invalid token")
}

// 仅管理员访问的中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := ParseToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
			c.Abort()
			return
		}

		// 验证角色是否为管理员
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have admin privileges"})
			c.Abort()
			return
		}

		// 将用户信息传递给后续处理程序
		c.Set("user", user)
		c.Next()
	}
}
