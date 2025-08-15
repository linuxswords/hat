package models

import "time"

// Score represents an archer's performance in a tournament
type Score struct {
	ID            int       `json:"id" db:"id"`
	ArcherID      int       `json:"archer_id" db:"archer_id"`
	TournamentID  int       `json:"tournament_id" db:"tournament_id"`
	BowClassID    string    `json:"bow_class_id" db:"bow_class_id"`
	RawScore      int       `json:"raw_score" db:"raw_score"`                     // The actual score shot
	AdjustedScore *float64  `json:"adjusted_score,omitempty" db:"adjusted_score"` // Raw score * handicap factor (optional)
	EnteredAt     time.Time `json:"entered_at" db:"entered_at"`
	EnteredBy     string    `json:"entered_by,omitempty" db:"entered_by"` // Who entered the score
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
