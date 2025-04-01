package controller

import (
	"rental-property-management-system/backend/config"
	"rental-property-management-system/backend/controller/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigRouter(r *gin.Engine) {
	// 代理静态文件
	// TODO: 最好移动到专门的网站上
	r.Static("/static", config.AppConfig.Server.StaticFilePath)

	// 无需特别权限的接口
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.GET("/list-rooms", ListRooms)
	r.POST("/get-room", GetRoom)
	r.POST("/list-filtered-rooms", ListFilteredRooms)

	// 需要身份认证的路由
	authRoutes := r.Group("/", middleware.AuthMiddleware())
	authRoutes.POST("/room", UpdateRoomInfo)
	authRoutes.POST("/create-order", CreateOrder)
	authRoutes.POST("/pay-bill", PayBill)
	authRoutes.GET("/list-relationships", ListRelationships)
	authRoutes.GET("/list-orders", ListOrders)
	authRoutes.GET("/list-billings", ListBillings)
	authRoutes.POST("/terminate-lease", TerminateLease)
	authRoutes.POST("/change-room-password", ChangeRoomPassword)
	authRoutes.POST("/get-password", GetPassword)
	authRoutes.GET("/list-user-workorders", ListUserWorkOrders)
	authRoutes.POST("/create-work-order", CreateWorkOrder)

	// 需要管理员权限的路由
	adminRoutes := authRoutes.Group("/", middleware.AdminMiddleware())
	adminRoutes.POST("/register-admin", RegisterAdmin)
	adminRoutes.POST("/update-room", UpdateRoomInfo)
	adminRoutes.POST("/admin-register", RegisterAdmin)
	adminRoutes.GET("/list-admin-workorders", ListAdminWorkOrders)
	adminRoutes.POST("/update-workorder", UpdateWorkOrder)
}
