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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// 类型转换
		user, ok := userAny.(store.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to convert any to store.User"})
			return
		}
		// 不是管理员身份无法通过当前中间件
		if user.Role != string(store.UserRoleAdmin) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
