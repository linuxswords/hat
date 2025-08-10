package models

// Handicap represents a handicap entry for a specific bow class
type Handicap struct {
	ID         int     `json:"id" db:"id"`
	BowClassID string  `json:"bow_class_id" db:"bow_class_id"`
	Factor     float64 `json:"factor" db:"factor"` // Handicap factor (e.g., 0.833, 1.0)
	SetID      int     `json:"set_id" db:"set_id"` // References HandicapSet
}

// HandicapSet represents a complete set of handicaps for all bow classes
type HandicapSet struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`        // e.g., "2024 World Championships"
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"` // Whether this set is currently in use
}
