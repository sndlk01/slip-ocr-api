package config

import (
	"log"
	"ocr-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	
	DB, err = gorm.Open(sqlite.Open(AppConfig.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Budget{},
		&models.Subscription{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database initialized successfully")
}

func GetDB() *gorm.DB {
	return DB
}
