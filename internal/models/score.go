package models

import "gorm.io/gorm"

type Score struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	TournamentID uint
	ArcherID     uint
	Score        float64
	TotalScore   float64
	Ranking      uint
}
