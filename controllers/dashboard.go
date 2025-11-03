package controllers

import (
	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/gorm"
	"strconv"
)

// DashboardController — faqat text yozuvlar bilan ishlaydi
type DashboardController struct {
	beego.Controller
}

// GET — dashboard sahifasini ko‘rsatish
func (c *DashboardController) Get() {
	idStr := c.Ctx.Input.Param(":id") // URL'dan id olish
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(400, "Noto'g'ri ID")
		return
	}

	var user models.User
	if err := database.DB.Preload("Notes", func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}).First(&user, userID).Error; err != nil {
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
