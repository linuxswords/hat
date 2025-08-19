package models

import "gorm.io/gorm"

// Handicap represents a handicap entry for a specific bow class
type Handicap struct {
	gorm.Model
	BowClassID string  `json:"bow_class_id" gorm:"not null;size:50;index" db:"bow_class_id"`
	Factor     float64 `json:"factor" gorm:"not null" db:"factor"` // Handicap factor (e.g., 0.833, 1.0)
	SetID      uint    `json:"set_id" gorm:"not null;index" db:"set_id"` // References HandicapSet
	
	// Foreign key relationship
	HandicapSet HandicapSet `json:"handicap_set,omitempty" gorm:"foreignKey:SetID"`
}

// HandicapSet represents a complete set of handicaps for all bow classes
type HandicapSet struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null;size:255" db:"name"`        // e.g., "2024 World Championships"
	Description string `json:"description" gorm:"size:500" db:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:false" db:"is_active"` // Whether this set is currently in use
	
	// Has many relationship
	Handicaps []Handicap `json:"handicaps,omitempty" gorm:"foreignKey:SetID"`
}
