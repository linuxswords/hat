package models

import (
	"gorm.io/gorm"
)

type TournamentArcher struct {
	gorm.Model
	ID              uint          `gorm:"primaryKey"`
	ArcherID        uint          `gorm:"not null"`
	Archer          Archer        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BowClassID      uint          `gorm:"not null"`
	BowClass        BowClass      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ScoreID         uint          `gorm:"not null"`
	Score           Score         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HandicapEntryID uint          `gorm:"not null"`
	HandicapEntry   HandicapEntry `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (ta *TournamentArcher) Name() string {
	if ta.Archer.ID == 0 {
		return "Unknown Archer"
	}
	return ta.Archer.Name()
}
