package models

import "time"

type Note struct {
	ID        uint   `gorm:"primaryKey"`
	Body      string `gorm:"type:text;not null"`
	Completed bool   `gorm:"default:false"`
	ImagePath string `gorm:"size:255"`
	CreatedAt time.Time

	// Foydalanuvchi bilan bog‘lanish
	UserID uint // FK (foreign key)
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
