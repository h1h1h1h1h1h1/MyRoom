package main

import (
	"database/sql"
	"fmt"
	"log"

	"electricity-management-backend/config"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 加载配置
	config.LoadConfig()
	cfg := config.AppConfig

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Database connection successful!")

	// 检查用户名是否已存在
	var count int
	username := "testuser3"
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check username: %v", err)
	}
	fmt.Printf("Username '%s' exists: %d\n", username, count)

	// 检查邮箱是否已存在
	email := "test3@example.com"
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check email: %v", err)
	}
	fmt.Printf("Email '%s' exists: %d\n", email, count)

	// 加密密码
