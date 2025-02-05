package database

import (
	"fmt"
	"io/ioutil"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rental-property-management-system/internal/config"
	"rental-property-management-system/internal/models"
)

var DB *gorm.DB

// ConnectDB 连接数据库
func ConnectDB() {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB = db
	initDatabase()
}

// initDatabase 执行init.sql文件初始化数据库
func initDatabase() error {
	// 读取init.sql文件内容
	sqlBytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql: %v", err)
	}

	// 执行SQL脚本
	// 使用DB.Exec执行原始SQL
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get raw SQL connection: %v", err)
	}

	_, err = sqlDB.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %v", err)
	}

	log.Println("Database initialized successfully.")
	return nil
}

// DisconnectDB 关闭数据库连接
func DisconnectDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database instance: " + err.Error())
	}
	err = sqlDB.Close()
	if err != nil {
		panic("Failed to close database connection: " + err.Error())
	}
}

// MigrateModels 自动迁移模型
func MigrateModels() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Password{},
	)
	if err != nil {
		panic("Database migration failed: " + err.Error())
	}
}
