package api

import (
	"database/sql"
	"net/http"

	"electricity-management-backend/database"
	"electricity-management-backend/models"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct{}

func NewInfoHandler() *InfoHandler {
	return &InfoHandler{}
}

func (h *InfoHandler) GetAnnouncements(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, type, summary, content, publish_date, expiry_date, status, is_important, address, phone, affected_area, created_at FROM information WHERE status = 'active' ORDER BY publish_date DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var infos []models.Information
	for rows.Next() {
		var info models.Information
		var address, phone, affectedArea sql.NullString
		// Note: publish_date and expiry_date in DB are DATE, scanning into time.Time works if parseTime=True
		if err := rows.Scan(&info.ID, &info.Title, &info.Type, &info.Summary, &info.Content, &info.PublishDate, &info.ExpiryDate, &info.Status, &info.IsImportant, &address, &phone, &affectedArea, &info.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan information"})
			return
		}
		info.Address = address.String
		info.Phone = phone.String
		info.AffectedArea = affectedArea.String
		infos = append(infos, info)
	}

	if infos == nil {
		infos = []models.Information{}
	}

	c.JSON(http.StatusOK, infos)
}
