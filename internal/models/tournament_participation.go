package models

import (
	"time"
	"gorm.io/gorm"
)

// TournamentParticipation represents an archer's participation in a tournament
type TournamentParticipation struct {
	gorm.Model
	TournamentID uint      `json:"tournament_id" gorm:"not null;index" db:"tournament_id"`
	ArcherID     uint      `json:"archer_id" gorm:"not null;index" db:"archer_id"`
	RegisteredAt time.Time `json:"registered_at" gorm:"autoCreateTime" db:"registered_at"`
	Status       string    `json:"status" gorm:"not null;size:50;default:'registered'" db:"status"` // "registered", "checked_in", "competing", "completed", "withdrawn"
	
	// Foreign key relationships
	Tournament Tournament `json:"tournament,omitempty" gorm:"foreignKey:TournamentID"`
	Archer     Archer     `json:"archer,omitempty" gorm:"foreignKey:ArcherID"`
}

// TournamentArcherView represents an archer with tournament-specific information for display
type TournamentArcherView struct {
	Participation TournamentParticipation `json:"participation"`
	Archer        Archer                  `json:"archer"`
	BowClassName  string                  `json:"bow_class_name,omitempty"`
}