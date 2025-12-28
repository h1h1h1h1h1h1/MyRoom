package api

import (
	"net/http"
	"time"

	"electricity-management-backend/database"
	"electricity-management-backend/middleware"
	"electricity-management-backend/models"
	"electricity-management-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// 检查邮箱是否已存在
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 插入用户
	result, err := database.DB.Exec(
		"INSERT INTO users (username, password, email, phone, real_name) VALUES (?, ?, ?, ?, ?)",
		req.Username, string(hashedPassword), req.Email, req.Phone, req.RealName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	userID, _ := result.LastInsertId()
	
	// 生成token
	token, err := utils.GenerateToken(int(userID), req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user": gin.H{
			"id":       userID,
			"username": req.Username,
			"email":    req.Email,
			"status":   "active",
		},
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, password, email FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 暂时移除状态检查
	// if user.Status != "active" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Account is not active"})
	// 	return
	// }

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"status":   "active", // 使用默认值
		},
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var user models.UserResponse
	err := database.DB.QueryRow(
		"SELECT id, username, email, phone, real_name, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.RealName, &user.CreatedAt)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 设置默认状态
	user.Status = "active"
	c.JSON(http.StatusOK, user)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建更新语句
	query := "UPDATE users SET updated_at = ?"
	params := []interface{}{time.Now()}

	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", req.Email, userID).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		}
		query += ", email = ?"
		params = append(params, req.Email)
	}

	if req.Phone != "" {
		query += ", phone = ?"
		params = append(params, req.Phone)
	}

	if req.RealName != "" {
		query += ", real_name = ?"
		params = append(params, req.RealName)
	}

	if req.Address != "" {
		query += ", address = ?"
		params = append(params, req.Address)
	}

	query += " WHERE id = ?"
	params = append(params, userID)

	_, err := database.DB.Exec(query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// GetCustomerNumbers 获取用户绑定的客户编号
func (h *UserHandler) GetCustomerNumbers(c *gin.Context) {
	userID := middleware.GetUserID(c)

	rows, err := database.DB.Query(
		"SELECT id, customer_number, meter_number, address, voltage_level, contract_capacity, created_at FROM customer_numbers WHERE user_id = ?",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var customers []models.CustomerNumber
	for rows.Next() {
		var customer models.CustomerNumber
		err := rows.Scan(
			&customer.ID, &customer.CustomerNumber, &customer.MeterNumber,
			&customer.Address, &customer.VoltageLevel, &customer.ContractCapacity,
			&customer.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data"})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, customers)
}

// BindCustomerNumber 绑定客户编号
func (h *UserHandler) BindCustomerNumber(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.BindCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查客户编号是否已被绑定
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM customer_numbers WHERE customer_number = ?", req.CustomerNumber).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer number already bound"})
		return
	}

	// 插入客户编号
	_, err = database.DB.Exec(
		"INSERT INTO customer_numbers (user_id, customer_number, meter_number, address) VALUES (?, ?, ?, ?)",
		userID, req.CustomerNumber, req.MeterNumber, req.Address,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind customer number"})
		return
	}

	// 初始化余额记录
	_, err = database.DB.Exec(
		"INSERT INTO electricity_balance (customer_number_id, balance) VALUES ((SELECT id FROM customer_numbers WHERE customer_number = ?), 0)",
		req.CustomerNumber,
	)
	if err != nil {
		// 如果余额初始化失败，记录错误但不影响主要操作
		c.JSON(http.StatusCreated, gin.H{
			"message": "Customer number bound successfully (balance initialization failed)",
			"warning": "Balance record not created",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer number bound successfully"})
}
