package middleware

import (
	"net/http"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

// 仅管理员访问的中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get(GinContextKeyUser)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// 类型转换
		user, ok := userAny.(store.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert any to store.User"})
			c.Abort()
			return
		}
		// 不是管理员身份无法通过当前中间件
		if user.Role != string(store.UserRoleAdmin) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
