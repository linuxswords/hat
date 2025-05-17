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

func (hs *HandycapSet) GetHandycapEntryByBowClass(bowClassID uint) *HandycapEntry {
	//FIXME this will retrieve the first match and will disregard other values from different handycap sets
	for _, entry := range hs.HandycapEntries {
		if entry.BowClassID == bowClassID {
			return &entry
		}
	}
	return nil
}

type HandycapEntry struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	BowClassID    uint
	BowClass      BowClass
	Value         float64
	HandycapSetID uint
}
