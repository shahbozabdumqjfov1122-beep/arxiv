package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"strings"

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
	c.TplName = "Buyurtma.html"
}

func (c *RegisterController) BuyurtmaPost() {
	c.TplName = "Buyurtma.html"
}

// GET - ro‘yxatdan o‘tish formasi
func (c *RegisterController) Get() {
	c.TplName = "register.html"
}

// POST - foydalanuvchi ma’lumotini saqlash
func (c *RegisterController) Post() {
	name := c.GetString("name")
	email := c.GetString("email")
	password := c.GetString("password")

	// 1️⃣ Bo‘sh maydonlarni tekshirish
	if name == "" || email == "" || password == "" {
		c.Data["Error"] = "⚠️ Barcha maydonlarni to‘ldiring!"
		c.TplName = "register.html"
		return
	}

	// 2️⃣ Emailni tekshirish
	if strings.Contains(email, "@gmail.com") {
		c.Data["Error"] = "❌ Iltimos, haqiqiy email kiriting (@gmail.com, @mail.ru va hokazo)."
		c.TplName = "register.html"
		return
	}

	if !strings.Contains(email, "@") {
		c.Data["Error"] = "❌ Email manzilda '@' belgisi bo‘lishi kerak."
		c.TplName = "register.html"
		return
	}

	// 3️⃣ Parolni hash qilish
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Data["Error"] = "❌ Parolni shifrlashda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// 4️⃣ User jadvaliga yozish
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.Data["Error"] = "❌ Foydalanuvchi saqlashda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// 5️⃣ Admin jadvaliga yozish (Role: User)
	admin := models.Admin{
		Firstname: name,
		Email:     email,
		Password:  string(hashed),
		Role:      "User",
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		c.Data["Error"] = "❌ Admin jadvaliga yozishda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// 6️⃣ Muvaffaqiyatli — login sahifasiga yo‘naltirish
	c.Redirect("/login", 302)
}
