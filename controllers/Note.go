package controllers

import (
	"os"
	"path/filepath"
	"strconv"

	"arxiv/database"
	"arxiv/models"
	beego "github.com/beego/beego/v2/server/web"
)

// NoteController — yozuvlarni o'chirish/tahrirlash uchun
type NoteController struct {
	beego.Controller
}

func (c *NoteController) Toggle() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.WriteString("Noto‘g‘ri ID")
		return
	}

	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.WriteString("Topilmadi")
		return
	}

	// Holatni o‘zgartirish (faqat bool bo‘lsa ishlaydi)
	note.Completed = !note.Completed
	database.DB.Save(&note)

	c.Ctx.Output.SetStatus(200)
	c.Ctx.WriteString("Tugallanish holati o‘zgartirildi")
}

// DELETE yoki GET orqali o'chirish
func (c *NoteController) Delete() {
	// Sessiyani tekshirish
	sessID := c.GetSession("user_id")
	if sessID == nil {
		c.CustomAbort(401, "Avtorizatsiya talab qilinadi")
		return
	}
	userID := sessID.(uint)

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.WriteString("Noto'g'ri ID")
		return
	}

	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.WriteString("Yozuv topilmadi")
		return
	}

	// Foydalanuvchi faqat o'z yozuvini o'chira oladi
	if note.UserID != userID {
		c.Ctx.Output.SetStatus(403)
		c.Ctx.WriteString("Siz bu yozuvni o'chira olmaysiz")
		return
	}

	// Agar rasm bo'lsa, faylni ham o'chirishga harakat qilamiz
	if note.ImagePath != "" {
		// pathni serverdagi haqiqiy joyga maplash
		if err := os.Remove(filepath.FromSlash(note.ImagePath)); err != nil {
			// Faylni o'chirishda xatolik bo'lsa ham davom etamiz
			// log qilish mumkin
		}
	}

	if err := database.DB.Delete(&models.Note{}, id).Error; err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Bazadan o'chirishda xatolik")
		return
	}

	// Redirect qayta dashboardga
	c.Redirect("/dashboard", 302)
}
