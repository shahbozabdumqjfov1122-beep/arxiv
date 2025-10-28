package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
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

	// passwordni hash qilish
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.WriteString("Parolni saqlashda xatolik")
		return
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPass),
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
