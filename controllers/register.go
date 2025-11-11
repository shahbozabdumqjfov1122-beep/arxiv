package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"fmt"
	"math/rand"
	"strings"
	"time"

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
	if email == "" || strings.HasSuffix(email, "@example.com") {
		c.Data["Error"] = "❌ Iltimos, haqiqiy email kiriting (masalan: @gmail.com, @mail.ru)."
		c.TplName = "register.html"
		return
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		c.Data["Error"] = "❌ Email manzilda '@' va domen bo‘lishi kerak (masalan: gmail.com)."
		c.TplName = "register.html"
		return
	}

	// 3️⃣ Parol uzunligini tekshirish
	if len(password) < 6 {
		c.Data["Error"] = "⚠️ Parol kamida 6 ta belgidan iborat bo‘lishi kerak!"
		c.TplName = "register.html"
		return
	}

	// 4️⃣ Foydalanuvchi allaqachon email bilan mavjudligini tekshirish
	var existing models.User
	result := database.DB.Where("email = ?", email).First(&existing)
	if result.RowsAffected > 0 {
		c.Data["Error"] = "⚠️ Bu email bilan foydalanuvchi allaqachon ro‘yxatdan o‘tgan!"
		c.TplName = "register.html"
		return
	}

	// 5️⃣ Username yagona bo‘lishini ta’minlash
	username := strings.ToLower(strings.TrimSpace(name))
	for {
		var existingUser models.User
		if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err != nil {
			break // username mavjud emas
		}
		// Tasodifiy qo‘shimcha qo‘shish
		rand.Seed(time.Now().UnixNano())
		username = fmt.Sprintf("%s_%d", strings.ToLower(strings.TrimSpace(name)), rand.Intn(1000))
	}

	// 6️⃣ Parolni hash qilish
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Data["Error"] = "❌ Parolni shifrlashda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// 7️⃣ User jadvaliga yozish
	user := models.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: string(hashed),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.Data["Error"] = "❌ Foydalanuvchi saqlashda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	// 8️⃣ Admin jadvaliga yozish (Role: User)
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

	// 9️⃣ Muvaffaqiyatli ro‘yxatdan o‘tganidan so‘ng
	c.Redirect("/login", 302)
}
