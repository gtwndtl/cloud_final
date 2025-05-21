package entity

import "gorm.io/gorm"

type Candidates struct {
	gorm.Model
	Name       string `json:"name"`
	ElectionID uint   `json:"election_id"`
}
