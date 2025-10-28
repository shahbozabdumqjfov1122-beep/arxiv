package controllers

import (
	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
)

// DashboardController — faqat text yozuvlar bilan ishlaydi
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
	// Notes’larni ID bo‘yicha kamayish tartibida olamiz
	database.DB.Preload("Notes", func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}).First(&user, userID)

	c.Data["User"] = user
	c.Data["UserId"] = user.ID
	c.TplName = "dashboard.html"
}

// POST — yangi note qo‘shish (faqat text)
// POST — yangi note qo‘shish (text yoki rasm)
func (c *DashboardController) Post() {
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.CustomAbort(401, "Avtorizatsiya talab qilinadi")
		return
	}
	userID := sessID.(uint)

	body := c.GetString("about")

	// 📸 Faylni olish (_ bilan, file ishlatilmaydi)
	_, header, err := c.GetFile("image")
	var imagePath string
	if err == nil && header != nil {
		// Fayl nomini yaratish
		imagePath = "static/uploads/" + header.Filename
		// Faylni saqlash
		if err := c.SaveToFile("image", imagePath); err != nil {
			c.Ctx.WriteString("Rasmni saqlashda xatolik: " + err.Error())
			return
		}
	}

	// Matn ham yo‘q, rasm ham yo‘q bo‘lsa, xato
	if body == "" && imagePath == "" {
		c.Data["json"] = map[string]interface{}{
			"error": "Hech qanday ma'lumot yuborilmadi",
		}
		c.ServeJSON()
		return
	}

	// 📝 Bazaga yozish
	note := models.Note{
		UserID:    userID,
		Body:      body,
		ImagePath: imagePath,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "Bazaga yozishda xatolik: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// ✅ JSON javob
	c.Data["json"] = map[string]interface{}{
		"ID":        note.ID,
		"Body":      note.Body,
		"ImagePath": note.ImagePath,
	}
	c.ServeJSON()
}
