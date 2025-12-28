package models

import "time"

type Application struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Status    string    `json:"status"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApplicationRequest struct {
	UserID  int    `json:"user_id"`
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type Information struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"` // service_point, outage_notice, news, policy
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	PublishDate  time.Time `json:"publish_date"`
	ExpiryDate   time.Time `json:"expiry_date"`
	Status       string    `json:"status"`
	IsImportant  bool      `json:"is_important"`
	Address      string    `json:"address,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	AffectedArea string    `json:"affected_area,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
