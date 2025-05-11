package models

import (
	"gorm.io/gorm"
)

type BowClass struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"unique;not null"`
	Name string `gorm:"not null"`
}
