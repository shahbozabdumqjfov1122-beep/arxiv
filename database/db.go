package database

import (
	"arxiv/models"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
)

var DB *gorm.DB

func InitDB() {
	runmode := beego.AppConfig.DefaultString("runmode", "dev")

	dbHost, _ := beego.AppConfig.String(runmode + "::db_host")
	dbPort, _ := beego.AppConfig.String(runmode + "::db_port")
	dbUser, _ := beego.AppConfig.String(runmode + "::db_user")
	dbPass, _ := beego.AppConfig.String(runmode + "::db_pass")
	dbName, _ := beego.AppConfig.String(runmode + "::db_name")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tashkent",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ DBga ulanish xatolik: %v", err)
	}
	DB = db

	// Migrate: avval Admin, keyin User va Note
	if err := DB.AutoMigrate(&models.Admin{}, &models.User{}, &models.Note{}); err != nil {
		log.Fatalf("❌ Migration xatolik: %v", err)
	}

	// Avtomatik admin yaratish
	SeedUserAdmin()
	log.Println("✅ DB tayyor va Admin seeding bajarildi")
}
func migrate() {
	err := DB.AutoMigrate(
		&models.Admin{}, // avval
		&models.User{},
		&models.Note{}, // oxir
	)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
