package controllers

import (
	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Help() {
	c.TplName = "help.html"
}

func (c *RegisterController) HelpPost() {
	c.TplName = "help.html"
}
func (c *RegisterController) Buyurtma() {
	c.TplName = "help.html"
}

func (c *RegisterController) BuyurtmaPost() {
	c.TplName = "help.html"
}

// GET - ro'yxatdan o'tish formasi
func (c *RegisterController) Get() {
	c.TplName = "register.html"
}

// POST - foydalanuvchi ma’lumotini saqlash
func (c *RegisterController) Post() {
	name := c.GetString("name")
	email := c.GetString("email")
	password := c.GetString("password")

	if name == "" || email == "" || password == "" {
		c.Data["Error"] = "Barcha maydonlarni to'ldiring!"
		c.TplName = "register.html"
		return
	}

	// Parolni hash qilish
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Data["Error"] = "Parolni shifrlashda xatolik"
		c.TplName = "register.html"
		return
	}

	// 1️⃣ User jadvaliga yozish
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.Data["Error"] = "Foydalanuvchi saqlashda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// 2️⃣ Admin jadvaliga yozish (Role: "User")
	admin := models.Admin{
		Firstname: name,
		Email:     email,
		Password:  string(hashed),
		Role:      "User",
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		c.Data["Error"] = "Admin jadvaliga yozishda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// muvaffaqiyatli saqlangandan keyin login sahifasiga yuborish
	c.Redirect("/login", 302)
}
