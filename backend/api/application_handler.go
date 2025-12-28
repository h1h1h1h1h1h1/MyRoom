package api

import (
	"net/http"

	"electricity-management-backend/database"
	"electricity-management-backend/middleware"
	"electricity-management-backend/models"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct{}

func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{}
}

func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	rows, err := database.DB.Query("SELECT id, user_id, type, content, status, created_at FROM applications WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var apps []models.Application
	for rows.Next() {
		var app models.Application
		if err := rows.Scan(&app.ID, &app.UserID, &app.Type, &app.Content, &app.Status, &app.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan application"})
			return
		}
		apps = append(apps, app)
	}

	if apps == nil {
		apps = []models.Application{}
	}

	c.JSON(http.StatusOK, apps)
}

func (h *ApplicationHandler) SubmitApplication(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req models.ApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("INSERT INTO applications (user_id, type, content, status) VALUES (?, ?, ?, 'pending')", userID, req.Type, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Application submitted successfully"})
}
