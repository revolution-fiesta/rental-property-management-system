package controller

import (
	"rental-property-management-system/backend/controller/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigRouter(r *gin.Engine) {
	// 用户注册和登录不需要管理员权限
	r.POST("/login", Login)
	r.POST("/register", Register)

	// 需要身份认证的路由
	authRoutes := r.Group("/", middleware.AuthMiddleware())
	authRoutes.POST("/room", UpdateRoomInfo)

	// 需要管理员权限的路由
	adminRoutes := authRoutes.Group("/", middleware.AdminMiddleware())
	adminRoutes.POST("/register-admin", RegisterAdmin)
	adminRoutes.POST("/update-room", UpdateRoomInfo)
	adminRoutes.POST("/admin-register", RegisterAdmin)
	adminRoutes.GET("/get-workorder", GetWorkOrdersByAdmin)
	adminRoutes.POST("/update-workorder", UpdateWorkOrderStatus)
}
