package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
