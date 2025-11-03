package models

import (
	"time"
)

type Note struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    uint   `gorm:"not null"`
	Body      string `gorm:"type:text;not null"`
	ImagePath string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
