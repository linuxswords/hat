package models

import (
	"gorm.io/gorm"
)

type HandycapSet struct {
	gorm.Model
	ID              uint `gorm:"primaryKey"`
	Name            string
	HandycapEntries []HandycapEntry `gorm:"foreignKey:HandycapSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type HandycapEntry struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	BowClassID    uint
	BowClass      BowClass
	Value         float64
	HandycapSetID uint
}
