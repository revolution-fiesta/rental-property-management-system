package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rental-property-management-system/backend/config"
	"rental-property-management-system/backend/controller"
	"rental-property-management-system/backend/runner"
	"rental-property-management-system/backend/utils"
	"runtime"

	"rental-property-management-system/backend/store"
	"syscall"

	"github.com/gin-gonic/gin"
)

// WARN: 超级大 bug: Find 函数没有找到结果的时候不会返回错误, 需要用 affected_rows 判断
func main() {
	// 用于性能测试
	runtime.GOMAXPROCS(8)

	var err error
	// 读取配置文件
	if err = config.LoadConfig(); err != nil {
		slog.Error(err.Error())
		return
	}

	// 初始化存储层, 并设置进程退出时关闭存储层
	if err = store.Init(); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		postgreErr, redisErr := store.Close()
		if redisErr != nil {
			slog.Warn(postgreErr.Error())
		}
		if postgreErr != nil {
			slog.Warn(redisErr.Error())
		}
	}()

	// 迁移数据库模型
	if err = store.Migrate(); err != nil {
		slog.Error(err.Error())
		return
	}

	// 生成服务器验证的公私钥对
	config.PrivateKey, err = utils.GeneratePrivateKey()
	if err != nil {
		slog.Error("failed to generate rsa key pairs")
		return
	}

	// 配置路由
	ginHandler := gin.Default()
	controller.ConfigRouter(ginHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.AppConfig.Server.Port),
		Handler: ginHandler,
	}

	// 处理 SIGTERM 以及 SIGINT 信号
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		slog.Info(fmt.Sprintf("%s received, bye~", sig.String()))
		cancel()
	}()

	// 启动服务器
	go func() {
		slog.Info(fmt.Sprintf("Server runs on port %s", config.AppConfig.Server.Port))
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Server encounted an error")
		}
		cancel()
	}()

	// 每隔一定时间检查一次数据库的账单进行出账
	runner.StartBillingRunner(ctx)
	// 阻塞直到 ctx 结束
	<-ctx.Done()
}
