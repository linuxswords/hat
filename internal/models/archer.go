package models

// Archer represents an athlete who participates in archery competitions/tournaments
type Archer struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Gender   string `json:"gender" db:"gender"`
	BowClass string `json:"bow_class" db:"bow_class"`
	Email    string `json:"email" db:"email"`
}
