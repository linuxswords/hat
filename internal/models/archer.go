package models

import (
	"gorm.io/gorm"
)

type Archer struct {
	gorm.Model
	BowClassID uint
	BowClass   BowClass
	FirstName string   `gorm:"not null"`
	LastName  string   `gorm:"not null"`
	Gender    string   `gorm:"not null"`
}
