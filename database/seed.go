package database

import (
	"arxiv/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func SeedUserAdmin() {
	var admin models.Admin

	if err := DB.Where("role = ?", "Admin").First(&admin).Error; err != nil {
		password := "123"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		newAdmin := models.Admin{
			Firstname: "Admin",
			Email:     "admin@example.com",
			Role:      "Admin",
			Password:  string(hashed),
		}

		if err := DB.Create(&newAdmin).Error; err != nil {
			log.Printf("❌ Admin yaratishda xatolik: %v", err)
			return
		}

		log.Println("✅ Admin foydalanuvchi muvaffaqiyatli yaratildi.")
	} else {
		log.Println("ℹ️ Admin allaqachon mavjud.")
	}
}
