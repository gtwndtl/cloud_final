package config

import (
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=election_db user=postgres password=postgres dbname=election_db port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	}

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
	// สร้างตาราง Elections
	err := db.AutoMigrate(&entity.Elections{})
	if err != nil {
		log.Fatal("Failed to migrate Elections table:", err)
	}

	// Seed ข้อมูลตัวอย่าง (optional)
	election := entity.Elections{
		Title:       "University Student Council 2025",
		Description: "เลือกตั้งสภานักศึกษา ประจำปีการศึกษา 2568",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(48 * time.Hour),
		Status:      "upcoming",
	}

	// สร้างเฉพาะเมื่อยังไม่มีข้อมูลนี้
	db.FirstOrCreate(&election, entity.Elections{Title: "University Student Council 2025"})
}
