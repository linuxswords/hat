package models

import (
	"gorm.io/gorm"
	"time"
)

type Tournament struct {
	gorm.Model
	ID                uint `gorm:"primaryKey"`
	Name              string
	Date              time.Time
	Venue             string
	TournamentType    string      `gorm:"not null"`
	HandicapSet       HandicapSet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HandicapSetID     uint
	TournamentArchers []TournamentArcher `gorm:"foreignKey:TournamentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
