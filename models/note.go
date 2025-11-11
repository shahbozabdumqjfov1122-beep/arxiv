package models

import "time"

type Note struct {
	ID        uint   `gorm:"primaryKey"`
	Body      string `gorm:"type:text;not null"`
	Completed bool   `gorm:"default:false"`
	ImagePath string `gorm:"size:255"`
	Name      string `gorm:"size:100"`
	Username  string `gorm:"size:100"`
	Email     string `gorm:"size:100"` // ✅ unique olib tashlandi
	CreatedAt time.Time

	// Foydalanuvchi bilan bog‘lanish
	UserID uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
