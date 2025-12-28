package main

import (
	"database/sql"
	"fmt"
	"log"

	"electricity-management-backend/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	var createTable string
	var tableName string
	err = db.QueryRow("SHOW CREATE TABLE users").Scan(&tableName, &createTable)
	if err != nil {
		log.Printf("Failed to get create table users: %v", err)
	} else {
		fmt.Printf("Users table definition:\n%s\n", createTable)
	}
}
