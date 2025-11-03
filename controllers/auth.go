package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	beego.Controller
}

// Login sahifasini ko'rsatish
func (c *AuthController) GetLogin() {
	c.TplName = "login.html"
}

// Login POST
func (c *AuthController) PostLogin() {
	email := strings.TrimSpace(c.GetString("email"))
	pass := strings.TrimSpace(c.GetString("password"))

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.Data["Error"] = "Email yoki parol noto'g'ri"
		c.TplName = "login.html"
		return
	}

	// parol tekshirish
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)) != nil {
		c.Data["Error"] = "Email yoki parol noto'g'ri"
		c.TplName = "login.html"
		return
	}

	// sessiyaga saqlash
	c.SetSession("user_id", user.ID)
	c.Redirect("/dashboard", 302) // ID ni URL'dan olib tashlang
}

// Register sahifasi
func (c *AuthController) GetRegister() {
	c.TplName = "register.html"
}

// Register POST
func (c *AuthController) PostRegister() {
	email := strings.TrimSpace(c.GetString("email"))
	pass := strings.TrimSpace(c.GetString("password"))
	name := strings.TrimSpace(c.GetString("name"))

	// email allaqachon mavjudligini tekshirish
	var exist models.User
	if err := database.DB.Where("email = ?", email).First(&exist).Error; err == nil {
		c.Data["Error"] = "Bu email allaqachon ro'yxatdan o'tgan"
		c.TplName = "register.html"
		return
	}

	// parolni hash qilish
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		c.Data["Error"] = "Parolni saqlashda xatolik"
		c.TplName = "register.html"
		return
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPass),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.Data["Error"] = "Ro'yxatdan o'tishda xatolik: " + err.Error()
		c.TplName = "register.html"
		return
	}

	c.Redirect("/login", 302)
}

// Logout
func (c *AuthController) Logout() {
	c.DelSession("user_id")
	c.Redirect("/login", 302)
}
