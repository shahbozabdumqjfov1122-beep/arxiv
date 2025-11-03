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

	// 📸 Faylni olish
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

	// Matn ham yo‘q, rasm ham yo‘q bo‘lsa
	if body == "" && imagePath == "" {
		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
		return
	}

	// 📝 Bazaga yozish
	note := models.Note{
		UserID:    userID,
		Body:      body,
		ImagePath: imagePath,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
		return
	}

	// ✅ Sahifani qayta yuklash
	c.Redirect("/dashboard", 302)
}

//package controllers
//
//import (
//	"arxiv/database"
//	"arxiv/models"
//	beego "github.com/beego/beego/v2/server/web"
//	"gorm.io/gorm"
//)
//
//// DashboardController — faqat text yozuvlar bilan ishlaydi
//type DashboardController struct {
//	beego.Controller
//}
//
//// GET — dashboard sahifasini ko‘rsatish
//func (c *DashboardController) Get() {
//	sessID := c.GetSession("user_id")
//	if sessID == nil {
//		c.Redirect("/login", 302)
//		return
//	}
//	userID := sessID.(uint)
//
//	var user models.User
//	// Notes’larni ID bo‘yicha kamayish tartibida olamiz
//	database.DB.Preload("Notes", func(db *gorm.DB) *gorm.DB {
//		return db.Order("id DESC")
//	}).First(&user, userID)
//
//	c.Data["User"] = user
//	c.Data["UserId"] = user.ID
//
//	// === YANGI QO'SHILGAN QISM ===
//	// URL da 'success' parametri borligini tekshiramiz
//	if c.GetString("success") == "1" {
//		c.Data["success"] = true // Agar bor bo'lsa, templatega 'success' ni yuboramiz
//	}
//	// === QO'SHILGAN QISM TUGADI ===
//
//	c.TplName = "dashboard.html"
//}
//
//func (c *DashboardController) Post() {
//	sessID := c.GetSession("user_id")
//	if sessID == nil {
//		c.CustomAbort(401, "Avtorizatsiya talab qilinadi")
//		return
//	}
//	userID := sessID.(uint)
//
//	body := c.GetString("about")
//
//	// 📸 Faylni olish
//	_, header, err := c.GetFile("image")
//	var imagePath string
//	if err == nil && header != nil {
//		// Fayl nomini yaratish
//		imagePath = "static/uploads/" + header.Filename
//		// Faylni saqlash
//		if err := c.SaveToFile("image", imagePath); err != nil {
//			c.Ctx.WriteString("Rasmni saqlashda xatolik: " + err.Error())
//			return
//		}
//	}
//
//	// Matn ham yo‘q, rasm ham yo‘q bo‘lsa
//	if body == "" && imagePath == "" {
//		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
//		return
//	}
//
//	// 📝 Bazaga yozish
//	note := models.Note{
//		UserID:    userID,
//		Body:      body,
//		ImagePath: imagePath,
//	}
//	if err := database.DB.Create(&note).Error; err != nil {
//		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
//		return
//	}
//}
