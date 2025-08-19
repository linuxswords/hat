package models

import "gorm.io/gorm"

// Archer represents an athlete who participates in archery competitions/tournaments
type Archer struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null;size:255" db:"name"`
	Gender   string `json:"gender" gorm:"not null;size:10" db:"gender"`
	BowClass string `json:"bow_class" gorm:"not null;size:50" db:"bow_class"`
	Email    string `json:"email" gorm:"not null;unique;size:255" db:"email"`
}
