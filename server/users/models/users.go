package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Email       string         `json:"email" gorm:"unique"`
	PhoneNumber string         `json:"-"`
	Password    string         `json:"-"`
	IsVerified  bool           `json:"is_verified"`
	Profile     Profile
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
