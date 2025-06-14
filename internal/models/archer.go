package models

import (
	"gorm.io/gorm"
)

type Archer struct {
	gorm.Model
	ID         uint     `gorm:"primaryKey"`
	BowClassID uint     `gorm:"not null"`
	BowClass   BowClass `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FirstName  string   `gorm:"not null"`
	LastName   string   `gorm:"not null"`
}

func (a *Archer) Name() string {
	return a.FirstName + " " + a.LastName
}
