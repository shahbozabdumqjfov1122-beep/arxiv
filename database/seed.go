package database

import (
	"arxiv/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func SeedUserAdmin() {
	var admin models.Admin

	// Admin mavjudligini tekshiramiz
	if err := DB.Where("role = ?", "Admin").First(&admin).Error; err != nil {
		// Yangi admin yaratamiz
		password := "123"

		// Parolni hash qilamiz
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		newAdmin := models.Admin{
			Firstname: "Admin",             // ism
			Email:     "admin@example.com", // login uchun email
			Role:      "Admin",
			Password:  string(hashed), // hashed parol
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
