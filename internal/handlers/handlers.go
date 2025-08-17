package handlers

import (
	"github.com/linuxswords/hat/internal/models"
)

// Repository interfaces
type ArcherRepository interface {
	GetAll() []models.Archer
	GetByID(id int) (*models.Archer, error)
	Create(archer models.Archer) (*models.Archer, error)
	Update(id int, archer models.Archer) (*models.Archer, error)
	Delete(id int) error
	GetByBowClass(bowClass string) []models.Archer
	GetByGender(gender string) []models.Archer
}

type TournamentRepository interface {
	GetAll() []models.Tournament
	GetByID(id int) (*models.Tournament, error)
	Create(tournament models.Tournament) (*models.Tournament, error)
	Update(id int, tournament models.Tournament) (*models.Tournament, error)
	Delete(id int) error
	GetUpcoming() []models.Tournament
	GetCurrent() []models.Tournament
	GetByHandicapSet(handicapSetID int) []models.Tournament
}

type HandicapRepository interface {
	GetAllSets() []models.HandicapSet
	GetSetByID(id int) (*models.HandicapSet, error)
	GetHandicapsBySetID(setID int) []models.Handicap
	GetHandicapByBowClass(setID int, bowClassID string) (*models.Handicap, error)
	GetActiveSet() (*models.HandicapSet, error)
}

type ScoreRepository interface {
	GetAll() []models.Score
	GetByTournamentID(tournamentID int) []models.Score
	GetByTournamentIDSorted(tournamentID int) []models.Score
	GetByArcherID(archerID int) []models.Score
	GetByTournamentAndArcher(tournamentID int, archerID int) (*models.Score, error)
	Create(score models.Score) (*models.Score, error)
	Update(id int, score models.Score) (*models.Score, error)
	UpdateByTournamentAndArcher(tournamentID int, archerID int, score models.Score) (*models.Score, error)
	Delete(id int) error
	DeleteByTournamentAndArcher(tournamentID int, archerID int) error
	GetScoreCountByTournament(tournamentID int) int
	CalculateAdjustedScore(rawScore int, handicapFactor float64) float64
}

type TournamentParticipationRepository interface {
	GetAll() []models.TournamentParticipation
	GetByTournamentID(tournamentID int) []models.TournamentParticipation
	GetByArcherID(archerID int) []models.TournamentParticipation
	GetByTournamentAndArcher(tournamentID int, archerID int) (*models.TournamentParticipation, error)
	Create(participation models.TournamentParticipation) (*models.TournamentParticipation, error)
	Update(id int, participation models.TournamentParticipation) (*models.TournamentParticipation, error)
	Delete(id int) error
	DeleteByTournamentAndArcher(tournamentID int, archerID int) error
	GetArcherCountByTournament(tournamentID int) int
	GetTournamentCountByArcher(archerID int) int
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
