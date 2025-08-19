package handlers

import (
	"github.com/linuxswords/hat/internal/models"
)

// Repository interfaces
type ArcherRepository interface {
	GetAll() []models.Archer
	GetByID(id uint) (*models.Archer, error)
	Create(archer models.Archer) (*models.Archer, error)
	Update(id uint, archer models.Archer) (*models.Archer, error)
	Delete(id uint) error
	GetByBowClass(bowClass string) []models.Archer
	GetByGender(gender string) []models.Archer
}

type TournamentRepository interface {
	GetAll() []models.Tournament
	GetByID(id uint) (*models.Tournament, error)
	Create(tournament models.Tournament) (*models.Tournament, error)
	Update(id uint, tournament models.Tournament) (*models.Tournament, error)
	Delete(id uint) error
	GetUpcoming() []models.Tournament
	GetCurrent() []models.Tournament
	GetByHandicapSet(handicapSetID uint) []models.Tournament
}

type HandicapRepository interface {
	GetAllSets() []models.HandicapSet
	GetSetByID(id uint) (*models.HandicapSet, error)
	GetHandicapsBySetID(setID uint) []models.Handicap
	GetHandicapByBowClass(setID uint, bowClassID string) (*models.Handicap, error)
	GetActiveSet() (*models.HandicapSet, error)
	CalculateAdjustedScore(rawScore int, handicapFactor float64) float64
}

type ScoreRepository interface {
	GetAll() []models.Score
	GetByTournamentID(tournamentID uint) []models.Score
	GetByTournamentIDSorted(tournamentID uint) []models.Score
	GetByArcherID(archerID uint) []models.Score
	GetByTournamentAndArcher(tournamentID uint, archerID uint) (*models.Score, error)
	Create(score models.Score) (*models.Score, error)
	Update(id uint, score models.Score) (*models.Score, error)
	UpdateByTournamentAndArcher(tournamentID uint, archerID uint, score models.Score) (*models.Score, error)
	Delete(id uint) error
	DeleteByTournamentAndArcher(tournamentID uint, archerID uint) error
	GetScoreCountByTournament(tournamentID uint) int64
	CalculateAdjustedScore(rawScore int, handicapFactor float64) float64
}

type TournamentParticipationRepository interface {
	GetAll() []models.TournamentParticipation
	GetByTournamentID(tournamentID uint) []models.TournamentParticipation
	GetByArcherID(archerID uint) []models.TournamentParticipation
	GetByTournamentAndArcher(tournamentID uint, archerID uint) (*models.TournamentParticipation, error)
	Create(participation models.TournamentParticipation) (*models.TournamentParticipation, error)
	Update(id uint, participation models.TournamentParticipation) (*models.TournamentParticipation, error)
	Delete(id uint) error
	DeleteByTournamentAndArcher(tournamentID uint, archerID uint) error
	GetArcherCountByTournament(tournamentID uint) int64
	GetTournamentCountByArcher(archerID uint) int64
}

// Handler structs with injected dependencies
type ArcherHandlers struct {
	ArcherRepo ArcherRepository
}

type TournamentHandlers struct {
	TournamentRepo TournamentRepository
	HandicapRepo   HandicapRepository
}

type TournamentScoreHandlers struct {
	ScoreRepo        ScoreRepository
	TournamentRepo   TournamentRepository
	ArcherRepo       ArcherRepository
	HandicapRepo     HandicapRepository
	ParticipationRepo TournamentParticipationRepository
}

type TournamentArcherHandlers struct {
	TournamentRepo    TournamentRepository
	ArcherRepo        ArcherRepository
	ParticipationRepo TournamentParticipationRepository
}

type HandicapHandlers struct {
	HandicapRepo HandicapRepository
}
