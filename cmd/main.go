package main

import (
	"rental-property-management-system/internal/config"
	"rental-property-management-system/internal/controllers"
	"rental-property-management-system/internal/store"

	//"rental-property-management-system/internal/models"
	"rental-property-management-system/api"
	"rental-property-management-system/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化配置
	config.LoadConfig()

	// 初始化数据库
	store.Init()

	// 在程序结束时关闭数据库连接
	defer store.Close()

	// 迁移数据库模型
	store.MigrateModels()
	// 初始化房间数据
	controllers.InitRoomData()

	// 创建Gin实例
	r := gin.Default()

	// 启动服务器
	r.POST("login", api.Login)
	r.POST("/register", api.Register)
	// 在路由中使用 AdminRequired 中间件，确保只有管理员可以访问
	r.PUT("/room", middleware.AdminRequired(), controllers.UpdateRoomInfo)
	r.Run(":8080")

}
