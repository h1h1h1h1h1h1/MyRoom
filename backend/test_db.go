package main

import (
	"database/sql"
	"fmt"
	"log"

	"electricity-management-backend/config"
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

	// 检查users表是否存在
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'users'", cfg.DBName).Scan(&tableExists)
	if err != nil {
		log.Fatalf("Failed to check table existence: %v", err)
	}

	if tableExists > 0 {
		fmt.Println("Users table exists!")
		
		// 检查表结构
		rows, err := db.Query("DESCRIBE users")
		if err != nil {
			log.Fatalf("Failed to describe users table: %v", err)
		}
		defer rows.Close()
		
		fmt.Println("\nUsers table structure:")
		for rows.Next() {
			var field, typ, null, key, extra string
			var defaultValue sql.NullString
			err := rows.Scan(&field, &typ, &null, &key, &defaultValue, &extra)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%-20s %-30s %-5s %-5s\n", field, typ, null, key)
		}
	} else {
		fmt.Println("Users table does not exist!")
		
		// 尝试创建表
		fmt.Println("\nCreating users table...")
		createSQL := `CREATE TABLE users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			phone VARCHAR(20),
			real_name VARCHAR(50),
			id_card VARCHAR(18),
			address VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			status ENUM('active', 'inactive', 'suspended') DEFAULT 'active'
		)`
		
		_, err := db.Exec(createSQL)
		if err != nil {
			log.Fatalf("Failed to create users table: %v", err)
		}
		fmt.Println("Users table created successfully!")
	}
}
