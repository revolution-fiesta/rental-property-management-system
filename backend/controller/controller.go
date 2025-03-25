package controller

import (
	"rental-property-management-system/backend/controller/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigRouter(r *gin.Engine) {
	// 用户注册和登录不需要管理员权限
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.PUT("/room", middleware.AdminRequired(), UpdateRoomInfo)

	// 需要管理员权限的路由
	r.PUT("/update-room", middleware.AdminRequired(), UpdateRoomInfo)
	r.PUT("/admin-register", middleware.AdminRequired(), RegisterAdmin)
	r.PUT("/get-workorder", middleware.AdminRequired(), GetWorkOrdersByAdmin)
	r.PUT("/update-workorder", middleware.AdminRequired(), UpdateWorkOrderStatus)
}
