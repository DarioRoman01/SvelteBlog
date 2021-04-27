package models

type Like struct {
	UserID uint `gorm:"primaryKey"`
	PostID uint `gorm:"primaryKey"`
	Value  uint16
}
