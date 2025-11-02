package models

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	Firstname string
	Email     string
	Password  string
	Role      string
	ImagePath string `gorm:"size:255"` // rasm yo‘li
}
