package entity

import (
	"time"

	"gorm.io/gorm"
)

type Votes struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	CandidateID uint      `json:"candidate_id"`
	ElectionID  uint      `json:"election_id"`
	Timestamp   time.Time `json:"timestamp"`
}
