package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100"` // <-- Shu qo'shiladi
	Username  string `gorm:"unique;size:100"`
	Password  string
	CreatedAt time.Time
	Notes     []Note
	Email     string `gorm:"unique;size:100"` // 👈 Email (auth.go va register.go uchun)
}
