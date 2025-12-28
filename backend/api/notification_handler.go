package api

import (
	"net/http"

	"electricity-management-backend/database"
	"electricity-management-backend/middleware"
	"electricity-management-backend/models"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	rows, err := database.DB.Query("SELECT id, user_id, title, content, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Title, &notif.Content, &notif.IsRead, &notif.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan notification"})
			return
		}
		notifs = append(notifs, notif)
	}

	if notifs == nil {
		notifs = []models.Notification{}
	}

	c.JSON(http.StatusOK, notifs)
}
