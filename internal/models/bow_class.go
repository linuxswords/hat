package models

// BowClass represents a classification of a bow and the archer who uses it
type BowClass struct {
	ID          string `json:"id" db:"id"`           // e.g., "AMLB" for Adult Male Longbow
	Name        string `json:"name" db:"name"`       // e.g., "Adult Male Longbow"
	Description string `json:"description" db:"description"`
	AgeGroup    string `json:"age_group" db:"age_group"`     // e.g., "Adult"
	Gender      string `json:"gender" db:"gender"`           // e.g., "Male"
	BowType     string `json:"bow_type" db:"bow_type"`       // e.g., "Longbow", "Compound", "Recurve"
}
