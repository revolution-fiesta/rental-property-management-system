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
	authRoutes.PUT("/room", UpdateRoomInfo)

	// 需要管理员权限的路由
	adminRoutes := authRoutes.Group("/", middleware.AdminMiddleware())
	adminRoutes.POST("/registerAdmin", RegisterAdmin)
	adminRoutes.PUT("/update-room", UpdateRoomInfo)
	adminRoutes.PUT("/admin-register", RegisterAdmin)
	adminRoutes.PUT("/get-workorder", GetWorkOrdersByAdmin)
	adminRoutes.PUT("/update-workorder", UpdateWorkOrderStatus)
}
