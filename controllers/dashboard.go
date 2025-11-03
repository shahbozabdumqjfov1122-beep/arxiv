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

	if c.GetString("success") == "1" {
		c.Data["success"] = true
	}

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

	_, header, err := c.GetFile("image")
	var imagePath string
	if err == nil && header != nil {
		imagePath = "static/uploads/" + header.Filename
		if err := c.SaveToFile("image", imagePath); err != nil {
			c.Ctx.WriteString("Rasmni saqlashda xatolik: " + err.Error())
			return
		}
	}

	if body == "" && imagePath == "" {
		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
		return
	}

	note := models.Note{
		UserID:    userID,
		Body:      body,
		ImagePath: imagePath,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
		return
	}

	c.Redirect("/dashboard?success=1", 302)
}
