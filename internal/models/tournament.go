package models

import "time"

// Tournament represents a competition where multiple archers compete
type Tournament struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Location     string    `json:"location" db:"location"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
	HandicapSetID int      `json:"handicap_set_id,omitempty" db:"handicap_set_id"` // Optional, for handicap tournaments (0 means no handicap)
}
