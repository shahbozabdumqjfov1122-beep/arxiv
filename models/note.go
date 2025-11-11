package models

import "time"

type Note struct {
	ID        uint   `gorm:"primaryKey"`
	Body      string `gorm:"type:text;not null"`
	Completed bool   `gorm:"default:false"`
	ImagePath string `gorm:"size:255"`
	Name      string `gorm:"size:100"` // <-- Shu qo'shiladi
	Username  string `gorm:"unique;size:100"`

	Email     string `gorm:"unique;size:100"` // 👈 Email (auth.go va register.go uchun)
	CreatedAt time.Time

	// Foydalanuvchi bilan bog‘lanish
	UserID uint // FK (foreign key)
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
