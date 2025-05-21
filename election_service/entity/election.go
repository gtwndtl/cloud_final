package entity

import( "gorm.io/gorm"
		"time"
)

type Elections struct {
	gorm.Model

	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`

	// เก็บ ID ของ candidate เป็น list หรือแค่ ID เดียว
	// เช่น ถ้ารายการหลายคน อาจเก็บเป็น []uint หรือ json string
	CandidateIDs []uint `gorm:"-" json:"candidate_ids"` // ไว้ให้โค้ดใช้งาน แต่ไม่แมปลงฐานข้อมูล
}
