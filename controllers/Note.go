package controllers

import (
	"arxiv/database"
	"arxiv/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type NoteController struct {
	beego.Controller
}

// Toggle note completed
func (c *NoteController) Toggle() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		c.Ctx.Output.SetStatus(404)
		return
	}

	note.Completed = !note.Completed
	database.DB.Save(&note)
	c.Ctx.Output.SetStatus(200)
}

// Delete note
func (c *NoteController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	if err := database.DB.Delete(&models.Note{}, id).Error; err != nil {
		c.Ctx.Output.SetStatus(404)
		return
	}

	c.Ctx.Output.SetStatus(200)
}
