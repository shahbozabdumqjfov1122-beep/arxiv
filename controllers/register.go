package controllers

import (
	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
)

type RegisterController struct {
	beego.Controller
}

// GET - formani ko'rsatish
func (c *RegisterController) Get() {
	c.TplName = "register.html"
}
func (c *RegisterController) Help() {
	c.TplName = "help.html"
}
func (c *RegisterController) Buyurtma() {
	c.TplName = "Buyurtma.html"
}

// POST - foydalanuvchi ma’lumotini saqlash
func (c *RegisterController) Post() {
	name := c.GetString("name")
	email := c.GetString("email")
	password := c.GetString("password")

	// ✅ Parolni hash qilmasdan saqlash
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password, // to'g'ridan-to'g'ri saqlanadi
	}

	// GORM orqali saqlash
	if err := database.DB.Create(&user).Error; err != nil {
		c.Data["Error"] = "Foydalanuvchi saqlashda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// muvaffaqiyatli saqlangandan keyin login sahifasiga yuborish
	c.Redirect("/login", 302)
}
