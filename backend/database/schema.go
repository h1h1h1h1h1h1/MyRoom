package database

import (
	"database/sql"
	"log"
)

func InitSchema(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			phone VARCHAR(20),
			real_name VARCHAR(50),
			id_card VARCHAR(20),
			address VARCHAR(255),
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS customer_numbers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT UNSIGNED NOT NULL,
			customer_number VARCHAR(50) NOT NULL UNIQUE,
			meter_number VARCHAR(50) NOT NULL,
			address VARCHAR(255) NOT NULL,
			voltage_level VARCHAR(20),
			contract_capacity DECIMAL(10, 2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS applications (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT UNSIGNED NOT NULL,
			type VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS information (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL,
			summary TEXT,
			content TEXT NOT NULL,
			publish_date DATE NOT NULL,
			expiry_date DATE NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			is_important BOOLEAN DEFAULT FALSE,
			address VARCHAR(255),
			phone VARCHAR(50),
			affected_area VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT UNSIGNED NOT NULL,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS electricity_balance (
			id INT AUTO_INCREMENT PRIMARY KEY,
			customer_number_id INT NOT NULL,
			balance DECIMAL(10, 2) DEFAULT 0.00,
			last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id)
		)`,
		`CREATE TABLE IF NOT EXISTS electricity_bills (
			id INT AUTO_INCREMENT PRIMARY KEY,
			customer_number_id INT NOT NULL,
			billing_month DATE NOT NULL,
			consumption DECIMAL(10, 2) NOT NULL,
			unit_price DECIMAL(10, 4) NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			due_date DATE NOT NULL,
			paid_amount DECIMAL(10, 2) DEFAULT 0.00,
			payment_status VARCHAR(20) DEFAULT 'unpaid',
			paid_date TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id)
		)`,
		`CREATE TABLE IF NOT EXISTS meter_readings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			customer_number_id INT NOT NULL,
			reading_date DATE NOT NULL,
			current_reading DECIMAL(10, 2) NOT NULL,
			previous_reading DECIMAL(10, 2),
			consumption DECIMAL(10, 2) NOT NULL,
			reading_type VARCHAR(20) DEFAULT 'manual',
			reader_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Failed to execute query: %s\nError: %v", query, err)
		}
	}
	log.Println("Database schema initialized successfully")
}
