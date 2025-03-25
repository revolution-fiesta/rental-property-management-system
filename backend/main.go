package main

import (
	"log/slog"
	"rental-property-management-system/backend/config"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置文件
	if err := config.LoadConfig(); err != nil {
		slog.Error(err.Error())
	}
	// 初始化存储层
	if err := store.Init(); err != nil {
		slog.Error(err.Error())
	}
	// 在程序结束时关闭数据库连接
	defer store.Close()
	// 迁移数据库模型
	if err := store.MigrateModels(); err != nil {
		slog.Error(err.Error())
	}
	// TODO: 仅用于测试
	if err := store.GenerateMockData(); err != nil {
		slog.Error(err.Error())
	}
	// 运行服务器
	gin.Default().Run(":8080")
}
