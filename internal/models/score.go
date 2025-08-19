package models

import (
	"time"
	"gorm.io/gorm"
)

// Score represents an archer's performance in a tournament
type Score struct {
	gorm.Model
	ArcherID      uint      `json:"archer_id" gorm:"not null;index" db:"archer_id"`
	TournamentID  uint      `json:"tournament_id" gorm:"not null;index" db:"tournament_id"`
	BowClassID    string    `json:"bow_class_id" gorm:"not null;size:50;index" db:"bow_class_id"`
	RawScore      int       `json:"raw_score" gorm:"not null" db:"raw_score"`                     // The actual score shot
	AdjustedScore *float64  `json:"adjusted_score,omitempty" db:"adjusted_score"` // Raw score * handicap factor (optional)
	EnteredAt     time.Time `json:"entered_at" gorm:"autoCreateTime" db:"entered_at"`
	EnteredBy     string    `json:"entered_by,omitempty" gorm:"size:255" db:"entered_by"` // Who entered the score
	
	// Foreign key relationships
	Archer       Archer     `json:"archer,omitempty" gorm:"foreignKey:ArcherID"`
	Tournament   Tournament `json:"tournament,omitempty" gorm:"foreignKey:TournamentID"`
}

// TournamentScoreView represents a score with archer and bow class information for display
type TournamentScoreView struct {
	Score                Score    `json:"score"`
	Archer               Archer   `json:"archer"`
	BowClassName         string   `json:"bow_class_name,omitempty"`
	HandicapFactor       float64  `json:"handicap_factor,omitempty"`
	Rank                 int      `json:"rank,omitempty"`
	DisplayAdjustedScore float64  `json:"display_adjusted_score,omitempty"` // Dereferenced adjusted score for templates
	HasAdjustedScore     bool     `json:"has_adjusted_score"`               // Whether adjusted score exists
}
