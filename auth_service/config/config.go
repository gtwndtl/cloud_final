package config

import (
	"log"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := "host=auth_db user=postgres password=postgres dbname=auth_db port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db = database
}


func DB() *gorm.DB {
	return db
}

func SetupDatabase() {
	// AutoMigrate สร้างตาราง
	err := db.AutoMigrate(&entity.Gender{}, &entity.Users{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed ข้อมูลเพศ
	genderMale := entity.Gender{Name: "Male"}
	genderFemale := entity.Gender{Name: "Female"}

	db.FirstOrCreate(&genderMale, entity.Gender{Name: "Male"})
	db.FirstOrCreate(&genderFemale, entity.Gender{Name: "Female"})

	// Seed ผู้ใช้งานตัวอย่าง
	hashedPassword := "123456" // ถ้ามีฟังก์ชัน HashPassword() ให้ใช้

	user := entity.Users{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Age:       30,
		Password:  hashedPassword, // ควรใช้ bcrypt หรือ hash จริงใน production
		Role:      "admin",
		BirthDay:  time.Date(1995, 5, 10, 0, 0, 0, 0, time.UTC),
		GenderID:  genderMale.ID,
	}

	db.FirstOrCreate(&user, entity.Users{Email: "admin@example.com"})
}