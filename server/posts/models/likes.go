package models

// like model
// the value only can be 1 or 0
type Like struct {
	UserID uint `gorm:"primaryKey"`
	PostID uint `gorm:"primaryKey"`
	Value  uint16
}
