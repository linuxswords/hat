package models

import (
	"time"
	"gorm.io/gorm"
)

// Tournament represents a competition where multiple archers compete
type Tournament struct {
	gorm.Model
	Name         string    `json:"name" gorm:"not null;size:255" db:"name"`
	Location     string    `json:"location" gorm:"not null;size:255" db:"location"`
	StartDate    time.Time `json:"start_date" gorm:"not null" db:"start_date"`
	EndDate      time.Time `json:"end_date" gorm:"not null" db:"end_date"`
	HandicapSetID uint     `json:"handicap_set_id,omitempty" gorm:"default:0" db:"handicap_set_id"` // Optional, for handicap tournaments (0 means no handicap)
}
