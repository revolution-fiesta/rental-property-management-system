package main

import (

	//"rental-property-management-system/internal/models"

	"log/slog"
	"rental-property-management-system/config"
	"rental-property-management-system/store"

	"github.com/gin-gonic/gin"
)

func handleErr(fn func() error) {
	if err := fn(); err != nil {
		slog.Error(err.Error())
	}
}

func main() {
	// 读取配置文件
	handleErr(config.LoadConfig)
	// 初始化存储层
	handleErr(store.Init)
	// 在程序结束时关闭数据库连接
	defer store.Close()
	// 迁移数据库模型
	handleErr(store.MigrateModels)
	// TODO: 仅用于测试
	handleErr(store.GenerateMockData)

	// 运行服务器
	gin.Default().Run(":8080")
}
