package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/se/entity"
)

var db *gorm.DB

func ConnectionDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// กำหนดค่า default ถ้าไม่ได้ตั้ง environment
		dsn = "postgresql://postgres:postgres@candidate_db:5432/candidate_db?sslmode=disable"
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
	err := db.AutoMigrate(&entity.Candidates{})
	if err != nil {
		log.Fatal("Failed to migrate candidate table:", err)
	}

	// Seed ตัวอย่าง
	candidate := entity.Candidates{
		Name:       "Alice Smith",
		ElectionID: 1, // ต้องรู้ว่ามี Election ID=1 มาจาก election_service
	}
	db.FirstOrCreate(&candidate, entity.Candidates{Name: "Alice Smith"})
}


