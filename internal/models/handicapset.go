package models

import (
	"gorm.io/gorm"
)

type HandycapSet struct {
	gorm.Model
	Name      string
	Handycaps []Handycap `gorm:"foreignKey:HandycapSetID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Handycap struct {
	gorm.Model
	BowClassID    uint
	BowClass      BowClass
	Value         float64
	HandycapSetID uint
}
