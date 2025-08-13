package models

import "time"

// TournamentParticipation represents an archer's participation in a tournament
type TournamentParticipation struct {
	ID           int       `json:"id" db:"id"`
	TournamentID int       `json:"tournament_id" db:"tournament_id"`
	ArcherID     int       `json:"archer_id" db:"archer_id"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
	Status       string    `json:"status" db:"status"` // "registered", "checked_in", "competing", "completed", "withdrawn"
}

// TournamentArcherView represents an archer with tournament-specific information for display
type TournamentArcherView struct {
	Participation TournamentParticipation `json:"participation"`
	Archer        Archer                  `json:"archer"`
	BowClassName  string                  `json:"bow_class_name,omitempty"`
}