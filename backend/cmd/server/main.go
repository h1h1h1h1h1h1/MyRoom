package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"electricity-management-backend/api"
	"electricity-management-backend/config"
	"electricity-management-backend/database"
	"electricity-management-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 初始化数据库表结构
	database.InitSchema(database.GetDB())

	// 设置Gin模式
	if config.AppConfig.ServerPort == "8080" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	router := gin.Default()

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080", "http://localhost:5174", "http://localhost:5175"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}))

	// 创建处理器
	userHandler := api.NewUserHandler()
	appHandler := api.NewApplicationHandler()
	infoHandler := api.NewInfoHandler()
	notifHandler := api.NewNotificationHandler()
	elecHandler := api.NewElectricityHandler()

	// 公开路由（无需认证）
	public := router.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// 需要认证的路由
	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		auth.GET("/profile", userHandler.GetProfile)
		auth.PUT("/profile", userHandler.UpdateProfile)
		auth.GET("/user/customers", userHandler.GetCustomerNumbers)
		auth.POST("/user/bind", userHandler.BindCustomerNumber)

		// 电费相关
		auth.GET("/electricity/usage", elecHandler.GetUsage)
		auth.GET("/electricity/payments", elecHandler.GetPayments)
		auth.POST("/electricity/pay", elecHandler.Pay)

		// 通知相关
		auth.GET("/user/notifications", notifHandler.GetNotifications)
		auth.PUT("/notifications/:id/read", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Mark notification as read endpoint"})
		})

		// 信息发布
		auth.GET("/info/announcements", infoHandler.GetAnnouncements)

		// 用电申请
		auth.GET("/user/applications", appHandler.GetApplications)
		auth.POST("/user/application", appHandler.SubmitApplication)

		// 系统检查
		auth.POST("/system/check", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// 管理路由（需要管理员权限）
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		// 管理员功能（待实现）
		admin.POST("/meter-readings", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Add meter reading endpoint"})
		})
		admin.POST("/generate-bills", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Generate bills endpoint"})
		})
		admin.POST("/notifications", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Send notification endpoint"})
		})
		admin.POST("/information", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Create information post endpoint"})
		})
	}

	// 启动服务器
	port := config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", port)

	// 优雅关闭
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
