package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"os"
	"path/filepath"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
)

type DashboardController struct {
	beego.Controller
}

// GET — dashboard sahifasini ko‘rsatish
func (c *DashboardController) Get() {
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.Redirect("/login", 302)
		return
	}
	userID := sessID.(uint)

	var user models.User
	err := database.DB.Preload("Notes", func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}).First(&user, userID).Error

	if err != nil {
		c.CustomAbort(404, "Foydalanuvchi topilmadi")
		return
	}

	c.Data["User"] = user
	c.Data["UserId"] = user.ID
	c.TplName = "dashboard.html"
}

// POST — yangi yozuv qo‘shish
func (c *DashboardController) Post() {
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.CustomAbort(401, "Avtorizatsiya talab qilinadi")
		return
	}
	userID := sessID.(uint)

	body := c.GetString("about")

	// Faylni olish
	file, header, err := c.GetFile("image")
	var imagePath string
	if err == nil && header != nil {
		defer file.Close()

		// upload papkasini yaratish (agar mavjud bo‘lmasa)
		uploadDir := "static/uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		// Fayl nomini xavfsiz shaklda yaratish
		filename := filepath.Base(header.Filename)
		imagePath = filepath.Join(uploadDir, filename)

		if err := c.SaveToFile("image", imagePath); err != nil {
			c.Ctx.WriteString("Rasmni saqlashda xatolik: " + err.Error())
			return
		}
	}

	// Matn ham, rasm ham yo‘q bo‘lsa
	if strings.TrimSpace(body) == "" && imagePath == "" {
		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
		return
	}

	// Yangi yozuv yaratish
	note := models.Note{
		UserID:    userID,
		Body:      body,
		ImagePath: imagePath,
	}

	if err := database.DB.Create(&note).Error; err != nil {
		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
		return
	}

	// Sahifani qayta yuklash
	c.Redirect("/dashboard", 302)
}
