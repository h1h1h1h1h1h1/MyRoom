package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
}

var AppConfig *Config

func LoadConfig() {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		port = 3306
	}

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "39.97.34.74"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "admin"),
		DBPassword: getEnv("DB_PASSWORD", "Jx200402!"),
		DBName:     getEnv("DB_NAME", "electricity_management"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "electricity-management-secret-key-2024"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	LoadConfig()
	log.Println("Configuration loaded successfully")
}
