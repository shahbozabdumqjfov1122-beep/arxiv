package models

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100"`
	Firstname string `gorm:"type:varchar(50);not null"`
	Email     string `gorm:"size:100;unique"`
	Password  string `gorm:"size:255"`
	Role      string `gorm:"type:varchar(50);not null"`
	Notes     []Note `gorm:"foreignKey:UserID"`
}
