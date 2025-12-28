package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	RealName  string    `json:"real_name,omitempty"`
	IDCard    string    `json:"id_card,omitempty"`
	Address   string    `json:"address,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone,omitempty"`
	RealName string `json:"real_name,omitempty"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	RealName string `json:"real_name,omitempty"`
	Address  string `json:"address,omitempty"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	RealName  string    `json:"real_name,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
