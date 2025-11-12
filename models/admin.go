package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin modeli - administratorlar uchun ma'lumotlar bazasi
type Admin struct {
	gorm.Model // GORM ning asosiy modeli (ID, CreatedAt, UpdatedAt, DeletedAt)

	// Asosiy ma'lumotlar
	Firstname string `gorm:"not null;size:100" json:"firstname"`          // Ism (majburiy)
	Email     string `gorm:"not null;unique;size:100" json:"email"`       // Email (majburiy, unikal)
	Password  string `gorm:"not null;size:255" json:"-"`                  // Hashlangan parol (majburiy)
	Role      string `gorm:"not null;default:'user';size:20" json:"role"` // Rol (majburiy, standart 'user')
	ImagePath string `gorm:"size:255" json:"image_path"`                  // Rasm yo'li (ixtiyoriy)
}

// Admin modelining metodlari
func (a *Admin) BeforeSave(tx *gorm.DB) error {
	// Parolni hash qilish (agar yangi parol kiritilgan bo'lsa)
	if len(a.Password) < 60 { // Agar parol hashlanmagan bo'lsa
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashedPass)
	}
	return nil
}
