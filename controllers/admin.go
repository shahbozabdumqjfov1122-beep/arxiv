package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AdminController struct {
	beego.Controller
}

// GET /admin — adminlar ro‘yxati
func (c *AdminController) Get() {
	var admins []models.Admin
	database.DB.Find(&admins)

	c.Data["Admins"] = admins
	c.TplName = "admin.html"
}

// POST /admin — yangi admin qo‘shish
func (c *AdminController) Post() {
	firstname := c.GetString("firstname")
	password := c.GetString("password")
	role := c.GetString("role")

	if firstname == "" || password == "" || role == "" {
		c.Ctx.WriteString("⚠️ Maydonlar to‘ldirilmagan!")
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	admin := models.Admin{
		Firstname: firstname,
		Password:  string(hashed),
		Role:      role,
	}

	database.DB.Create(&admin)
	c.Redirect("/admin", http.StatusFound)
}

// GET /admin/login — login sahifa
func (c *AdminController) Login() {
	c.TplName = "adminLogin.html"
}

// POST /admin/login — login tekshirish
func (c *AdminController) LoginPost() {
	email := strings.TrimSpace(c.GetString("email"))
	password := c.GetString("password")

	if email == "" || password == "" {
		c.Data["Error"] = "Iltimos, email va parolni kiriting."
		c.TplName = "adminLogin.html"
		return
	}

	var admin models.Admin
	if err := database.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		c.Data["Error"] = "Foydalanuvchi topilmadi."
		c.TplName = "adminLogin.html"
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		c.Data["Error"] = "Parol noto‘g‘ri."
		c.TplName = "adminLogin.html"
		return
	}

	// Sessiya yaratish
	c.SetSession("admin_id", admin.ID)
	c.Redirect("/admin", http.StatusFound)
}

// GET /admin/delete?id=1 — adminni o‘chirish
func (c *AdminController) Delete() {
	id, _ := c.GetInt("id")
	database.DB.Delete(&models.Admin{}, id)
	c.Redirect("/admin", http.StatusFound)
}
