package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
	"time"

	"arxiv/database"
	"arxiv/models"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

type DashboardController struct {
	beego.Controller
}

func (c *DashboardController) Get() {
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.Redirect("/login", 302)
		return
	}
	userID := sessID.(uint)

	var user models.User
	err := database.DB.Preload("Notes", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&user, userID).Error
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["User"] = user
	c.Data["Success"] = c.GetSession("success") != nil
	c.SetSession("success", nil)
	c.TplName = "dashboard.html"
}

func (c *DashboardController) Post() {
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.CustomAbort(401, "Avtorizatsiya talab qilinadi")
		return
	}
	userID := sessID.(uint)

	// Formani validatsiya qilish
	valid := validation.Validation{}
	body := c.GetString("about")
	valid.Required(body, "about").Message("Matn maydoni majburdir")

	// Faylni olish va validatsiya qilish
	file, fileHeader, err := c.GetFile("image")
	var imagePath string
	if err == nil && fileHeader != nil {
		defer file.Close()

		// Fayl formatini tekshirish
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		allowedExts := []string{".jpg", ".jpeg", ".png", ".gif"}
		isValidExt := false
		for _, allowed := range allowedExts {
			if ext == allowed {
				isValidExt = true
				break
			}
		}

		if !isValidExt {
			valid.SetError("image", "Fayl formati noto'g'ri. Faqat JPG, PNG va GIF fayllari qabul qilinadi")
		} else {
			// Fayl hajmini tekshirish (5MB gacha)
			if fileHeader.Size > 5*1024*1024 {
				valid.SetError("image", "Fayl hajmi 5MB dan oshmasligi kerak")
			} else {
				// Fayl nomini generatsiya qilish
				filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
				imagePath = filepath.Join("static/uploads", filename)

				// Faylni saqlash
				if err := c.SaveToFile("image", imagePath); err != nil {
					valid.SetError("image", "Rasmni saqlashda xatolik")
				}
			}
		}
	}

	// Validatsiya natijasini tekshirish
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.Ctx.WriteString(fmt.Sprintf("%s: %s\n", err.Field, err.Message))
		}
		return
	}

	// Matn ham, rasm ham bo'lmasa
	if body == "" && imagePath == "" {
		c.Ctx.WriteString("Hech qanday ma'lumot yuborilmadi")
		return
	}

	// Bazaga yozish
	note := models.Note{
		UserID:    userID,
		Body:      body,
		ImagePath: imagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := database.DB.Create(&note).Error; err != nil {
		c.Ctx.WriteString("Bazaga yozishda xatolik: " + err.Error())
		return
	}

	// Muvaffaqiyatli xabar
	c.SetSession("success", true)
	c.Redirect("/dashboard", 302)
}
