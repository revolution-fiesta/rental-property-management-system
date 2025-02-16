package main

import (
	"fmt"
	"rental-property-management-system/internal/config"
	"rental-property-management-system/internal/controllers"
	"rental-property-management-system/internal/database"
	//"rental-property-management-system/internal/models"
	"github.com/gin-gonic/gin"
	"rental-property-management-system/api"
	"rental-property-management-system/middleware"
)



func main() {
  
	// 初始化配置
	config.LoadConfig()

	// 初始化数据库
	database.ConnectDB()
	fmt.Println("success")
	defer database.DisconnectDB() // 在程序结束时关闭数据库连接

	// 迁移数据库模型
	database.MigrateModels()
	// 初始化房间数据
	controllers.InitRoomData()

	// 创建Gin实例
	r := gin.Default()

	// 启动服务器
	r.POST("login",api.Login)
	r.POST("/register", api.Register)
	// 在路由中使用 AdminRequired 中间件，确保只有管理员可以访问
	r.PUT("/room", middleware.AdminRequired(), controllers.UpdateRoomInfo)
	r.Run(":8080")
	
}

