package models

import (
	"gorm.io/gorm"
)

type BowClass struct {
	gorm.Model
	Code string `gorm:"unique;not null"`
	Name string `gorm:"not null"`
}
