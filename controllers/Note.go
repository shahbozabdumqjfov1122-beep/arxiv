package controllers

import (
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

// DELETE /note/delete/:id
func (c *NoteController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	// 🗑️ Yozuvni bazadan o‘chiramiz
	if err := database.DB.Delete(&models.Note{}, id).Error; err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.WriteString("O‘chirishda xatolik: " + err.Error())
		return
	}

	// ✅ O‘chirilgandan so‘ng status 200 qaytariladi
	c.Ctx.Output.SetStatus(200)
}
