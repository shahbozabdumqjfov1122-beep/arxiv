package controllers

import (
	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
	"time"
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

	// ✅ O'zbekiston vaqtiga o‘tkazish
	loc, _ := time.LoadLocation("Asia/Tashkent")
	for i := range user.Notes {
		user.Notes[i].CreatedAt = user.Notes[i].CreatedAt.In(loc)
	}

	c.Data["User"] = user
	c.Data["UserId"] = user.ID
	c.TplName = "dashboard.html"
}

func (c *DashboardController) Post() {
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.CustomAbort(401, "Avtorizatsiya talab qilinadi")
		return
	}
	userID := sessID.(uint)

	body := c.GetString("about")

	// Faylni olish
	_, header, err := c.GetFile("image")
	var imagePath string
	hasImage := (err == nil && header != nil)

	// Foydalanuvchi yozuvlari va rasmlari sonini tekshirish
	var totalNotes int64
	var totalImages int64

	database.DB.Model(&models.Note{}).Where("user_id = ?", userID).Count(&totalNotes)
	database.DB.Model(&models.Note{}).Where("user_id = ? AND image_path <> ''", userID).Count(&totalImages)

	// 🚫 Limitni tekshirish
	if !hasImage && totalNotes >= 200 {
		c.Data["LimitError"] = "❌ Siz 200 ta yozuvdan ortiq qo‘sha olmaysiz."
		c.Get()
		return
	}
	if hasImage && totalImages >= 30 {
		c.Data["LimitError"] = "❌ Siz 30 tadan kop rasm yuklay olmaysiz."
		c.Get()
		return
	}

	// Agar rasm bo‘lsa, saqlaymiz
	if hasImage {
		imagePath = "static/uploads/" + header.Filename
		if err := c.SaveToFile("image", imagePath); err != nil {
			c.Ctx.WriteString("Rasmni saqlashda xatolik: " + err.Error())
			return
		}
	}

	// Matn ham yo‘q, rasm ham yo‘q bo‘lsa
	if body == "" && imagePath == "" {
		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
		return
	}

	// Bazaga yozish
	note := models.Note{
		UserID:    userID,
		Body:      body,
		ImagePath: imagePath,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
		return
	}

	c.Redirect("/dashboard", 302)
}
