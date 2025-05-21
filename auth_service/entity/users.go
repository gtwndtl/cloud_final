package entity

import (
	"time"
)

type Users struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Age       uint8     `json:"age"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	BirthDay  time.Time `json:"birthday"`
	GenderID  uint      `json:"gender_id"`
	Gender    Gender    `gorm:"foreignKey:GenderID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Gender struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
