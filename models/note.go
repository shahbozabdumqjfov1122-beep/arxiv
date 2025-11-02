package models

type Note struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Body      string `gorm:"type:text"`
	ImagePath string `gorm:"size:255"`
	Completed bool
}
