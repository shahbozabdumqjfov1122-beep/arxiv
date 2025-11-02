package controllers

import (
	"arxiv/database"
	"arxiv/models"

	beego "github.com/beego/beego/v2/server/web"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Help() {
	c.TplName = "help.html"
}

func (c *RegisterController) Buyurtma() {
	c.TplName = "Buyurtma.html"
}

// GET /register — ro‘yxatdan o‘tish sahifasi
func (c *RegisterController) Get() {
	c.TplName = "register.html"
}

// POST /register — foydalanuvchini saqlash
func (c *RegisterController) Post() {
	name := c.GetString("name")
	email := c.GetString("email")
	password := c.GetString("password")

	// Yangi foydalanuvchi obyekt
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password, // ⚠️ ochiq parol
	}

	// Bazaga yozish
	if err := database.DB.Create(&user).Error; err != nil {
		c.Data["Error"] = "Xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	c.Redirect("/login", 302)
}
