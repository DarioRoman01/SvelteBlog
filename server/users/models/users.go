package models

import (
	"blogv2/posts/models"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey;foreignkey:UserID;references:ID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Email       string         `json:"email" gorm:"unique"`
	PhoneNumber string         `json:"-" gorm:"unique"`
	Password    string         `json:"-"`
	IsVerified  bool           `json:"is_verified" gorm:"default=false"`
	Profile     Profile
	Posts       []models.Post
	Comments    []models.Comment
}

type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterInput struct {
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}
