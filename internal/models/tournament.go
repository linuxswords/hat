package models

import (
	"gorm.io/gorm"
	"time"
)

type Tournament struct {
	gorm.Model
	Name          string
	Date          time.Time
	Venue         string
	HandycapSet   HandycapSet
	HandycapSetID uint
	Archers       []Archer `gorm:"many2many:tournament_archers;"`
}
