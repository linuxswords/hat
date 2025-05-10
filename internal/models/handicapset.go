package models

import (
	"gorm.io/gorm"
)

type HandicapSet struct {
	gorm.Model
	Name      string
	Handicaps []Handicap
}

type Handicap struct {
	gorm.Model
	BowClassID    uint
	BowClass      BowClass
	Value         float64
	HandicapSetID uint
	Name          string
}
