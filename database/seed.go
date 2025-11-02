package database

import (
	"arxiv/models"
	"log"
)

// SeedUserAdmin — agar bazada admin yo‘q bo‘lsa, avtomatik yaratadi
func SeedUserAdmin() {
	var admin models.Admin

	// Admin mavjudligini tekshiramiz
	if err := DB.Where("role = ?", "Admin").First(&admin).Error; err != nil {
		// Yangi admin yaratamiz
		newAdmin := models.Admin{
			Firstname: "Admin",
			Role:      "Admin",
			Password:  "123", // Istasangiz bu joyda hashing qo‘shing
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
