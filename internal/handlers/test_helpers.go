package handlers

import (
	"github.com/linuxswords/hat/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockArcherRepository implements ArcherRepository for testing
type MockArcherRepository struct {
	mock.Mock
}

func (m *MockArcherRepository) GetAll() []models.Archer {
	args := m.Called()
	return args.Get(0).([]models.Archer)
}

func (m *MockArcherRepository) GetByID(id uint) (*models.Archer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Archer), args.Error(1)
}

func (m *MockArcherRepository) Create(archer models.Archer) (*models.Archer, error) {
	args := m.Called(archer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Archer), args.Error(1)
}

func (m *MockArcherRepository) Update(id uint, archer models.Archer) (*models.Archer, error) {
	args := m.Called(id, archer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Archer), args.Error(1)
}

func (m *MockArcherRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockArcherRepository) GetByBowClass(bowClass string) []models.Archer {
	args := m.Called(bowClass)
	return args.Get(0).([]models.Archer)
}

func (m *MockArcherRepository) GetByGender(gender string) []models.Archer {
	args := m.Called(gender)
	return args.Get(0).([]models.Archer)
}

// MockTournamentRepository implements TournamentRepository for testing
type MockTournamentRepository struct {
	mock.Mock
}

func (m *MockTournamentRepository) GetAll() []models.Tournament {
	args := m.Called()
	return args.Get(0).([]models.Tournament)
}

func (m *MockTournamentRepository) GetByID(id uint) (*models.Tournament, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tournament), args.Error(1)
}

func (m *MockTournamentRepository) Create(tournament models.Tournament) (*models.Tournament, error) {
	args := m.Called(tournament)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tournament), args.Error(1)
}

func (m *MockTournamentRepository) Update(id uint, tournament models.Tournament) (*models.Tournament, error) {
	args := m.Called(id, tournament)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tournament), args.Error(1)
}

func (m *MockTournamentRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTournamentRepository) GetUpcoming() []models.Tournament {
	args := m.Called()
	return args.Get(0).([]models.Tournament)
}

func (m *MockTournamentRepository) GetCurrent() []models.Tournament {
	args := m.Called()
	return args.Get(0).([]models.Tournament)
}

func (m *MockTournamentRepository) GetByHandicapSet(handicapSetID uint) []models.Tournament {
	args := m.Called(handicapSetID)
	return args.Get(0).([]models.Tournament)
}

// MockHandicapRepository implements HandicapRepository for testing
type MockHandicapRepository struct {
	mock.Mock
}

func (m *MockHandicapRepository) GetAllSets() []models.HandicapSet {
	args := m.Called()
	return args.Get(0).([]models.HandicapSet)
}

func (m *MockHandicapRepository) GetSetByID(id uint) (*models.HandicapSet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.HandicapSet), args.Error(1)
}

func (m *MockHandicapRepository) GetHandicapsBySetID(setID uint) []models.Handicap {
	args := m.Called(setID)
	return args.Get(0).([]models.Handicap)
}

func (m *MockHandicapRepository) GetHandicapByBowClass(setID uint, bowClass string) (*models.Handicap, error) {
	args := m.Called(setID, bowClass)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Handicap), args.Error(1)
}

func (m *MockHandicapRepository) GetActiveSet() (*models.HandicapSet, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.HandicapSet), args.Error(1)
}

func (m *MockHandicapRepository) CalculateAdjustedScore(rawScore int, handicapFactor float64) float64 {
	args := m.Called(rawScore, handicapFactor)
	return args.Get(0).(float64)
}

// MockScoreRepository implements ScoreRepository for testing
type MockScoreRepository struct {
	mock.Mock
}

func (m *MockScoreRepository) GetAll() []models.Score {
	args := m.Called()
	return args.Get(0).([]models.Score)
}

func (m *MockScoreRepository) GetByTournamentID(tournamentID uint) []models.Score {
	args := m.Called(tournamentID)
	return args.Get(0).([]models.Score)
}

func (m *MockScoreRepository) GetByTournamentIDSorted(tournamentID uint) []models.Score {
	args := m.Called(tournamentID)
	return args.Get(0).([]models.Score)
}

func (m *MockScoreRepository) GetByArcherID(archerID uint) []models.Score {
	args := m.Called(archerID)
	return args.Get(0).([]models.Score)
}

func (m *MockScoreRepository) GetByTournamentAndArcher(tournamentID uint, archerID uint) (*models.Score, error) {
	args := m.Called(tournamentID, archerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Score), args.Error(1)
}

func (m *MockScoreRepository) Create(score models.Score) (*models.Score, error) {
	args := m.Called(score)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Score), args.Error(1)
}

func (m *MockScoreRepository) Update(id uint, score models.Score) (*models.Score, error) {
	args := m.Called(id, score)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Score), args.Error(1)
}

func (m *MockScoreRepository) UpdateByTournamentAndArcher(tournamentID uint, archerID uint, score models.Score) (*models.Score, error) {
	args := m.Called(tournamentID, archerID, score)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Score), args.Error(1)
}

func (m *MockScoreRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockScoreRepository) DeleteByTournamentAndArcher(tournamentID uint, archerID uint) error {
	args := m.Called(tournamentID, archerID)
	return args.Error(0)
}

func (m *MockScoreRepository) GetScoreCountByTournament(tournamentID uint) int64 {
	args := m.Called(tournamentID)
	return args.Get(0).(int64)
}

func (m *MockScoreRepository) CalculateAdjustedScore(rawScore int, handicapFactor float64) float64 {
	args := m.Called(rawScore, handicapFactor)
	return args.Get(0).(float64)
}

// MockTournamentParticipationRepository implements TournamentParticipationRepository for testing
type MockTournamentParticipationRepository struct {
	mock.Mock
}

func (m *MockTournamentParticipationRepository) GetAll() []models.TournamentParticipation {
	args := m.Called()
	return args.Get(0).([]models.TournamentParticipation)
}

func (m *MockTournamentParticipationRepository) GetByTournamentID(tournamentID uint) []models.TournamentParticipation {
	args := m.Called(tournamentID)
	return args.Get(0).([]models.TournamentParticipation)
}

func (m *MockTournamentParticipationRepository) GetByArcherID(archerID uint) []models.TournamentParticipation {
	args := m.Called(archerID)
	return args.Get(0).([]models.TournamentParticipation)
}

func (m *MockTournamentParticipationRepository) GetByTournamentAndArcher(tournamentID uint, archerID uint) (*models.TournamentParticipation, error) {
	args := m.Called(tournamentID, archerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TournamentParticipation), args.Error(1)
}

func (m *MockTournamentParticipationRepository) Create(participation models.TournamentParticipation) (*models.TournamentParticipation, error) {
	args := m.Called(participation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TournamentParticipation), args.Error(1)
}

func (m *MockTournamentParticipationRepository) Update(id uint, participation models.TournamentParticipation) (*models.TournamentParticipation, error) {
	args := m.Called(id, participation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TournamentParticipation), args.Error(1)
}

func (m *MockTournamentParticipationRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTournamentParticipationRepository) DeleteByTournamentAndArcher(tournamentID uint, archerID uint) error {
	args := m.Called(tournamentID, archerID)
	return args.Error(0)
}

func (m *MockTournamentParticipationRepository) GetArcherCountByTournament(tournamentID uint) int64 {
	args := m.Called(tournamentID)
	return args.Get(0).(int64)
}

func (m *MockTournamentParticipationRepository) GetTournamentCountByArcher(archerID uint) int64 {
	args := m.Called(archerID)
	return args.Get(0).(int64)
}