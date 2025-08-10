package models

// Score represents an archer's performance in a tournament
type Score struct {
	ID           int     `json:"id" db:"id"`
	ArcherID     int     `json:"archer_id" db:"archer_id"`
	TournamentID int     `json:"tournament_id" db:"tournament_id"`
	BowClassID   string  `json:"bow_class_id" db:"bow_class_id"`
	RawScore     int     `json:"raw_score" db:"raw_score"`         // The actual score shot
	AdjustedScore *float64 `json:"adjusted_score,omitempty" db:"adjusted_score"` // Raw score * handicap factor (optional)
}
