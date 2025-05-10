package models

import (
	"gorm.io/gorm"
)

type HandycapSet struct {
	gorm.Model
	Name      string
	Handicaps []Handycap `gorm:"foreignKey:HandicapSetID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Handycap struct {
	gorm.Model
	BowClassID    uint
	BowClass      BowClass
	Value         float64
	HandicapSetID uint
	Name          string
}
