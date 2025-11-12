package models

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	Firstname string
	Email     string `gorm:"unique;size:100"` // Unikal email cheklovi qo'shildi
	Password  string
	Role      string
	ImagePath string `gorm:"size:255"` // rasm yo‘li
}
