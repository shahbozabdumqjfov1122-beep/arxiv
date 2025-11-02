package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// AdminController — adminlar bilan ishlash uchun controller
type AdminController struct {
	beego.Controller
}

// GET /admin
// Adminlar ro‘yxatini ko‘rsatadi
func (c *AdminController) Get() {
	var admins []models.Admin
	database.DB.Find(&admins)

	c.Data["Admins"] = admins
	c.TplName = "admin.html"
}

// POST /admin/create
// Yangi admin yaratadi
func (c *AdminController) Post() {
	firstname := c.GetString("firstname")
	password := c.GetString("password")
	role := c.GetString("role")

	if firstname == "" || password == "" || role == "" {
		c.Ctx.WriteString("⚠️ Maydonlar to‘ldirilmagan!")
		return
	}

	// Parolni hash qilish
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Ctx.WriteString("❌ Parolni shifrlashda xatolik")
		return
	}

	admin := models.Admin{
		Firstname: firstname,
		Password:  string(hashed),
		Role:      role,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		c.Ctx.WriteString("❌ Admin yaratilmadi: " + err.Error())
		return
	}

	c.Redirect("/admin", http.StatusFound)
}

// GET /admin/delete?id=1
// Adminni o‘chirish
func (c *AdminController) Delete() {
	id, _ := c.GetInt("id")

	if err := database.DB.Delete(&models.Admin{}, id).Error; err != nil {
		c.Ctx.WriteString("❌ O‘chirishda xatolik: " + err.Error())
		return
	}

	c.Redirect("/admin", http.StatusFound)
}
