package models

import (
	"gorm.io/gorm"
)

type HandicapSet struct {
	gorm.Model
	ID              uint `gorm:"primaryKey"`
	Name            string
	HandicapEntries []HandicapEntry `gorm:"foreignKey:HandicapSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (hs *HandicapSet) GetHandicapEntryByBowClass(bowClassID uint) *HandicapEntry {
	//FIXME this will retrieve the first match and will disregard other values from different handicap sets
	for _, entry := range hs.HandicapEntries {
		if entry.BowClassID == bowClassID {
			return &entry
		}
	}
	return nil
}

type HandicapEntry struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	BowClassID    uint
	BowClass      BowClass
	Value         float64
	HandicapSetID uint
}
