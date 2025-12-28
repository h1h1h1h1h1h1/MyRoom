-- 电费管理系统数据库设计
-- MySQL数据库，IP: 39.97.34.74

-- 创建数据库
CREATE DATABASE IF NOT EXISTS electricity_management;
USE electricity_management;

-- 用户表
CREATE TABLE users (
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
);

-- 客户编号表（一个用户可以绑定多个客户编号）
CREATE TABLE customer_numbers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    customer_number VARCHAR(50) UNIQUE NOT NULL,
    meter_number VARCHAR(50) NOT NULL,
    address VARCHAR(255) NOT NULL,
    voltage_level VARCHAR(20),
    contract_capacity DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_customer_number (customer_number)
);

-- 电费余额表
CREATE TABLE electricity_balance (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_number_id INT NOT NULL,
    balance DECIMAL(10,2) DEFAULT 0.00,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE CASCADE,
    INDEX idx_customer_number_id (customer_number_id)
);

-- 电表读数表
CREATE TABLE meter_readings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_number_id INT NOT NULL,
    reading_date DATE NOT NULL,
    current_reading DECIMAL(10,2) NOT NULL,
    previous_reading DECIMAL(10,2),
    consumption DECIMAL(10,2) GENERATED ALWAYS AS (current_reading - previous_reading) STORED,
    reading_type ENUM('monthly', 'special', 'initial') DEFAULT 'monthly',
    reader_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE CASCADE,
    UNIQUE KEY uk_reading_date (customer_number_id, reading_date),
    INDEX idx_reading_date (reading_date)
);

-- 电费账单表
CREATE TABLE electricity_bills (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_number_id INT NOT NULL,
    billing_month DATE NOT NULL,
    consumption DECIMAL(10,2) NOT NULL,
    unit_price DECIMAL(6,4) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    due_date DATE NOT NULL,
    paid_amount DECIMAL(10,2) DEFAULT 0.00,
    payment_status ENUM('unpaid', 'partial', 'paid', 'overdue') DEFAULT 'unpaid',
    paid_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE CASCADE,
    UNIQUE KEY uk_billing_month (customer_number_id, billing_month),
    INDEX idx_billing_month (billing_month),
    INDEX idx_payment_status (payment_status)
);

-- 购电记录表
CREATE TABLE purchase_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_number_id INT NOT NULL,
    purchase_amount DECIMAL(10,2) NOT NULL,
    payment_method ENUM('alipay', 'wechat', 'bank_card', 'cash') NOT NULL,
    transaction_id VARCHAR(100) UNIQUE,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'completed', 'failed', 'refunded') DEFAULT 'completed',
    notes TEXT,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE CASCADE,
    INDEX idx_purchase_date (purchase_date),
    INDEX idx_customer_purchase (customer_number_id, purchase_date)
);

-- 通知表
CREATE TABLE notifications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    notification_type ENUM('balance_low', 'payment_due', 'power_outage', 'info', 'warning') NOT NULL,
    send_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_status BOOLEAN DEFAULT FALSE,
    read_date TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_notification (user_id, send_date),
    INDEX idx_read_status (read_status)
);

-- 余额不足通知记录表
CREATE TABLE balance_notifications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_number_id INT NOT NULL,
    notification_id INT NOT NULL,
    notification_type ENUM('first', 'second') NOT NULL,
    sent_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    next_notification_date DATE,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE CASCADE,
    FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE,
    INDEX idx_customer_notification (customer_number_id, sent_date)
);

-- 停电记录表
CREATE TABLE power_outages (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_number_id INT NOT NULL,
    outage_date DATE NOT NULL,
    outage_reason ENUM('non_payment', 'maintenance', 'emergency', 'other') NOT NULL,
    restoration_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE CASCADE,
    INDEX idx_outage_date (outage_date)
);

-- 信息发布表
CREATE TABLE information_posts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    post_type ENUM('service_points', 'outage_announcement', 'policy', 'news', 'other') NOT NULL,
    publish_date DATE NOT NULL,
    expiry_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT,
    INDEX idx_post_type (post_type),
    INDEX idx_publish_date (publish_date),
    INDEX idx_is_active (is_active)
);

-- 用电申请表
CREATE TABLE electricity_applications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    application_type ENUM('new_installation', 'name_change', 'meter_verification', 'other') NOT NULL,
    customer_number_id INT,
    application_date DATE NOT NULL,
    status ENUM('pending', 'processing', 'approved', 'rejected', 'completed') DEFAULT 'pending',
    description TEXT,
    required_documents TEXT,
    processing_notes TEXT,
    completed_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_number_id) REFERENCES customer_numbers(id) ON DELETE SET NULL,
    INDEX idx_application_type (application_type),
    INDEX idx_status (status),
    INDEX idx_application_date (application_date)
);

-- 服务网点表
CREATE TABLE service_points (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    business_hours VARCHAR(100),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_is_active (is_active)
);

-- 电价表
CREATE TABLE electricity_prices (
    id INT PRIMARY KEY AUTO_INCREMENT,
    voltage_level VARCHAR(20) NOT NULL,
    tier INT NOT NULL,
    min_consumption DECIMAL(10,2),
    max_consumption DECIMAL(10,2),
    unit_price DECIMAL(6,4) NOT NULL,
    effective_date DATE NOT NULL,
    expiry_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_price_tier (voltage_level, tier, effective_date),
    INDEX idx_effective_date (effective_date)
);

-- 插入初始数据：电价
INSERT INTO electricity_prices (voltage_level, tier, min_consumption, max_consumption, unit_price, effective_date) VALUES
('居民用电', 1, 0, 240, 0.4883, '2024-01-01'),
('居民用电', 2, 241, 400, 0.5383, '2024-01-01'),
('居民用电', 3, 401, NULL, 0.7883, '2024-01-01'),
('商业用电', 1, 0, NULL, 0.8563, '2024-01-01');

-- 插入初始数据：服务网点
INSERT INTO service_points (name, address, phone, business_hours, latitude, longitude) VALUES
('城东营业厅', 'XX市城东区人民路123号', '0571-88880001', '周一至周五 8:30-17:00', 30.123456, 120.123456),
('城南营业厅', 'XX市城南区中山路456号', '0571-88880002', '周一至周五 8:30-17:00', 30.223456, 120.223456),
('城西营业厅', 'XX市城西区解放路789号', '0571-88880003', '周一至周五 8:30-17:00', 30.323456, 120.323456);

-- 创建视图：用户用电概览
CREATE VIEW user_electricity_overview AS
SELECT 
    u.id as user_id,
    u.username,
    u.real_name,
    cn.customer_number,
    cn.address,
    eb.balance,
    (SELECT reading_date FROM meter_readings mr WHERE mr.customer_number_id = cn.id ORDER BY reading_date DESC LIMIT 1) as last_reading_date,
    (SELECT current_reading FROM meter_readings mr WHERE mr.customer_number_id = cn.id ORDER BY reading_date DESC LIMIT 1) as last_reading,
    (SELECT amount FROM electricity_bills ebill WHERE ebill.customer_number_id = cn.id AND ebill.payment_status != 'paid' ORDER BY billing_month DESC LIMIT 1) as latest_unpaid_amount
FROM users u
JOIN customer_numbers cn ON u.id = cn.user_id
LEFT JOIN electricity_balance eb ON cn.id = eb.customer_number_id;

-- 创建存储过程：生成月度电费账单
DELIMITER //
CREATE PROCEDURE generate_monthly_bills(IN billing_date DATE)
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE customer_id INT;
    DECLARE prev_reading DECIMAL(10,2);
    DECLARE curr_reading DECIMAL(10,2);
    DECLARE consumption DECIMAL(10,2);
    DECLARE unit_price DECIMAL(6,4);
    DECLARE amount DECIMAL(10,2);
    DECLARE due_date DATE;
    
    -- 游标：获取所有客户上个月的电表读数
    DECLARE cur CURSOR FOR 
    SELECT 
        mr.customer_number_id,
        mr.previous_reading,
        mr.current_reading
    FROM meter_readings mr
    WHERE mr.reading_date = DATE_SUB(billing_date, INTERVAL 1 MONTH)
      AND mr.reading_type = 'monthly';
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    OPEN cur;
    
    read_loop: LOOP
        FETCH cur INTO customer_id, prev_reading, curr_reading;
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        -- 计算用电量
        SET consumption = curr_reading - prev_reading;
        
        -- 获取电价（简化：使用居民用电第一阶梯）
        SELECT unit_price INTO unit_price 
        FROM electricity_prices 
        WHERE voltage_level = '居民用电' 
          AND tier = 1 
          AND effective_date <= billing_date
        ORDER BY effective_date DESC 
        LIMIT 1;
        
        -- 计算电费金额
        SET amount = consumption * unit_price;
        
        -- 设置缴费截止日期（当月20日）
        SET due_date = DATE_ADD(billing_date, INTERVAL 20 DAY);
        
        -- 插入电费账单
        INSERT INTO electricity_bills (customer_number_id, billing_month, consumption, unit_price, amount, due_date)
        VALUES (customer_id, billing_date, consumption, unit_price, amount, due_date)
        ON DUPLICATE KEY UPDATE 
            consumption = VALUES(consumption),
            unit_price = VALUES(unit_price),
            amount = VALUES(amount),
            due_date = VALUES(due_date);
    END LOOP;
    
    CLOSE cur;
END//
DELIMITER ;

-- 创建触发器：购电后更新余额
DELIMITER //
CREATE TRIGGER after_purchase_insert
AFTER INSERT ON purchase_records
FOR EACH ROW
BEGIN
    IF NEW.status = 'completed' THEN
        UPDATE electricity_balance 
        SET balance = balance + NEW.purchase_amount
        WHERE customer_number_id = NEW.customer_number_id;
        
        -- 如果记录不存在，则插入
        IF ROW_COUNT() = 0 THEN
            INSERT INTO electricity_balance (customer_number_id, balance)
            VALUES (NEW.customer_number_id, NEW.purchase_amount);
        END IF;
    END IF;
END//
DELIMITER ;

-- 创建触发器：支付电费后更新余额
DELIMITER //
CREATE TRIGGER after_bill_payment
AFTER UPDATE ON electricity_bills
FOR EACH ROW
BEGIN
    IF NEW.paid_amount > OLD.paid_amount THEN
        UPDATE electricity_balance 
        SET balance = balance - (NEW.paid_amount - OLD.paid_amount)
        WHERE customer_number_id = NEW.customer_number_id;
    END IF;
END//
DELIMITER ;
