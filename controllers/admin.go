package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"net/http"
	"os"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AdminController struct {
	beego.Controller
}

// GET /admin — adminlar ro‘yxati
func (c *AdminController) Get() {
	// ✅ Sessiya tekshirish
	adminID := c.GetSession("admin_id")
	if adminID == nil {
		// Agar sessiya yo‘q bo‘lsa, login sahifasiga yo‘naltir
		c.Redirect("/admin/login", 302)
		return
	}

	// Sessiya mavjud bo‘lsa, admin ro‘yxatini ko‘rsat
	var admins []models.Admin
	database.DB.Find(&admins)

	c.Data["Admins"] = admins
	c.TplName = "admin.html"
}
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

	// Sessiyani o‘rnatish
	c.SetSession("admin_id", admin.ID)
	c.SetSession("admin_email", admin.Email)

	c.Redirect("/admin", 302)
}
func (c *AdminController) Logout() {
	c.DestroySession()
	c.Redirect("/admin/login", 302)
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

// GET /admin/delete?id=1 — adminni o‘chirish
func (c *AdminController) Delete() {
	// 1️⃣ ID ni oling
	id, err := c.GetInt("id")
	if err != nil {
		c.Ctx.WriteString("❌ ID noto‘g‘ri")
		return
	}

	// 2️⃣ Adminni bazadan oling
	var admin models.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		c.Ctx.WriteString("❌ Admin topilmadi: " + err.Error())
		return
	}

	// 3️⃣ Agar ImagePath bo‘lsa, faylni diskdan o‘chirish
	if admin.ImagePath != "" {
		if err := os.Remove(admin.ImagePath); err != nil && !os.IsNotExist(err) {
			// Fayl topilmasa yoki boshqa xato bo‘lsa — xatolikni yozib qo‘yish
			c.Ctx.WriteString("⚠️ Image faylni o‘chirishda xatolik: " + err.Error())
			return
		}
	}

	// 4️⃣ Adminni bazadan o‘chirish
	if err := database.DB.Delete(&models.Admin{}, id).Error; err != nil {
		c.Ctx.WriteString("❌ Bazadan o‘chirishda xatolik: " + err.Error())
		return
	}

	// 5️⃣ Redirect
	c.Redirect("/admin", 302)
}
