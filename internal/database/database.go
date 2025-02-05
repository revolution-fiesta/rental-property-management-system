package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rental_system/internal/config"
	"rental_system/internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB = db
}

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