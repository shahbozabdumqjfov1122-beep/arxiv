package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"gorm.io/gorm"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// DashboardController — foydalanuvchi yozuvlari bilan ishlaydi
type DashboardController struct {
	beego.Controller
}

// GET — dashboard sahifasini ko‘rsatish
func (c *DashboardController) Get() {
	// URL’dan :id olish
	idStr := c.Ctx.Input.Param(":id")
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

	// User initial
	initial := ""
	if len(user.Username) > 0 {
		initial = string([]rune(user.Username)[0])
	}

	c.Data["User"] = user
	c.Data["UserId"] = user.ID
	c.Data["UserInitial"] = initial

	// URL parametrida 'success' borligini tekshirish
	if c.GetString("success") == "1" {
		c.Data["success"] = true
	}

	c.TplName = "dashboard.html"
}

// POST — yangi yozuv qo‘shish
func (c *DashboardController) Post() {
	idStr := c.Ctx.Input.Param(":id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(400, "Noto'g'ri ID")
		return
	}

	body := c.GetString("about")

	// Faylni olish
	_, header, err := c.GetFile("image")
	var imagePath string
	if err == nil && header != nil {
		imagePath = "static/uploads/" + header.Filename
		if err := c.SaveToFile("image", imagePath); err != nil {
			c.Ctx.WriteString("Rasmni saqlashda xatolik: " + err.Error())
			return
		}
	}

	// Hech narsa yuborilmasa
	if body == "" && imagePath == "" {
		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
		return
	}

	// Bazaga yozish
	note := models.Note{
		UserID:    uint(userID),
		Body:      body,
		ImagePath: imagePath,
	}
	if err := database.DB.Create(&note).Error; err != nil {
		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
		return
	}

	// Redirect + success
	c.Redirect("/dashboard/"+strconv.Itoa(userID)+"?success=1", 302)
}
