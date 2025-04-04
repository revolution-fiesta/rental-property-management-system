package middleware

import (
	"log/slog"
	"net/http"
	"rental-property-management-system/backend/config"
	"rental-property-management-system/backend/store"
	"rental-property-management-system/backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// 验证用户身份的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		// 用 [7:] 去掉 "Bearer " 前缀
		userId, err := utils.ValidateAccessToken(tokenString[7:], &config.PrivateKey.PublicKey)
		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// 根据 ID 从数据库查询用户
		var user store.User
		if err := store.GetDB().First(&user, userId).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 将用户 ID 传递给后续处理程序
		c.Set(GinContextKeyUser, user)
		c.Next()
	}
}
