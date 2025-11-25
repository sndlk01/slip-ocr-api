package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort      string
	DatabasePath    string
	UploadDir       string
	TesseractLang   string
	MaxUploadSize   int64 
}

var AppConfig *Config

func Init() {
	AppConfig = &Config{
		ServerPort:      getEnv("SERVER_PORT", "8077"),
		DatabasePath:    getEnv("DATABASE_PATH", "./db.sqlite"),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		TesseractLang:   getEnv("TESSERACT_LANG", "tha+eng"), // Thai + English
		MaxUploadSize:   10 * 1024 * 1024, // 10MB
	}

	if err := os.MkdirAll(AppConfig.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	log.Println("Configuration initialized successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
