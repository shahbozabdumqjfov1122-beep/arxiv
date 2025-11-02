package controllers

import (
	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AdminController struct {
	beego.Controller
}

// ... (avvalgi kod) ...

// GET /admin/login
func (c *AdminController) Login() {
	// Agar xatolik xabari bo'lsa, uni templatega yuboramiz
	c.TplName = "adminLogin.html"
}

// POST /admin/login
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
		// topilmadi
		c.Data["Error"] = "Foydalanuvchi topilmadi."
		c.TplName = "adminLogin.html"
		return
	}

	// Parolni tekshirish:
	// - agar admin.Password bcrypt hash bo'lsa, bcrypt bilan solishtiramiz
	// - aks holda oddiy matn solishtirish (ehtiyot bo'ling: buni aniqroq qilish uchun yaratishda parolni hash qiling)
	useBcrypt := false
	if len(admin.Password) > 3 && (admin.Password[0:4] == "$2a$" || admin.Password[0:4] == "$2b$" || admin.Password[0:4] == "$2y$") {
		useBcrypt = true
	}

	if useBcrypt {
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
			c.Data["Error"] = "Parol noto‘g‘ri."
			c.TplName = "adminLogin.html"
			return
		}
	} else {
		// Eslatma: agar parollar xom matn sifatida saqlangan bo'lsa — xavfsizlik zaifligi.
		if admin.Password != password {
			c.Data["Error"] = "Parol noto‘g‘ri."
			c.TplName = "adminLogin.html"
			return
		}
	}

	// Kirish muvaffaqiyatli — sessiya yoki cookie o'rnating (oddiy redirect ko'rsatish)
	// beego sessiya ishlatayotgan bo'lsangiz misol:
	c.SetSession("admin_id", admin.ID)
	c.SetSession("admin_email", admin.Email)
	c.Redirect("/admin", http.StatusFound)
}
