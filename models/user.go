package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100"`
	Email    string `gorm:"size:100;unique"`
	Password string `gorm:"size:255"`
	Notes    []Note `gorm:"foreignKey:UserID"`
	Username string
}
