package api

import (
	"rental-property-management-system/controllers"
	"rental-property-management-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 用户注册和登录不需要管理员权限
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.PUT("/room", middleware.AdminRequired(), controllers.UpdateRoomInfo)

	// 需要管理员权限的路由
	r.PUT("/update-room", middleware.AdminRequired(), controllers.UpdateRoomInfo)
	r.PUT("/admin-register", middleware.AdminRequired(), RegisterAdmin)
	r.PUT("/get-workorder", middleware.AdminRequired(), controllers.GetWorkOrdersByAdmin)
	r.PUT("/update-workorder", middleware.AdminRequired(), controllers.UpdateWorkOrderStatus)
}
