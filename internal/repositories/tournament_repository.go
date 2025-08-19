package repositories

import (
	"fmt"
	"time"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// TournamentRepo implements TournamentRepository using GORM
type TournamentRepo struct {
	db *gorm.DB
}

// NewTournamentRepository creates a new database-backed tournament repository
func NewTournamentRepository(db *gorm.DB) *TournamentRepo {
	return &TournamentRepo{db: db}
}

// GetAll returns all tournaments
func (r *TournamentRepo) GetAll() []models.Tournament {
	var tournaments []models.Tournament
	r.db.Order("start_date DESC").Find(&tournaments)
	return tournaments
}

// GetByID returns a tournament by ID
func (r *TournamentRepo) GetByID(id uint) (*models.Tournament, error) {
	var tournament models.Tournament
	result := r.db.First(&tournament, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tournament with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &tournament, nil
}

// Create creates a new tournament
func (r *TournamentRepo) Create(tournament models.Tournament) (*models.Tournament, error) {
	result := r.db.Create(&tournament)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tournament, nil
}

// Update updates an existing tournament
func (r *TournamentRepo) Update(id uint, updatedTournament models.Tournament) (*models.Tournament, error) {
	var tournament models.Tournament
	result := r.db.First(&tournament, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tournament with ID %d not found", id)
		}
		return nil, result.Error
	}
	
	// Update fields
	tournament.Name = updatedTournament.Name
	tournament.Location = updatedTournament.Location
	tournament.StartDate = updatedTournament.StartDate
	tournament.EndDate = updatedTournament.EndDate
	tournament.HandicapSetID = updatedTournament.HandicapSetID
	
	result = r.db.Save(&tournament)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &tournament, nil
}

// Delete deletes a tournament by ID
func (r *TournamentRepo) Delete(id uint) error {
	result := r.db.Delete(&models.Tournament{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tournament with ID %d not found", id)
	}
	return nil
}

// GetUpcoming returns tournaments that haven't started yet
func (r *TournamentRepo) GetUpcoming() []models.Tournament {
	var tournaments []models.Tournament
	now := time.Now()
	r.db.Where("start_date > ?", now).Order("start_date ASC").Find(&tournaments)
	return tournaments
}

// GetCurrent returns tournaments that are currently running
func (r *TournamentRepo) GetCurrent() []models.Tournament {
	var tournaments []models.Tournament
	now := time.Now()
	r.db.Where("start_date <= ? AND end_date >= ?", now, now).Order("start_date ASC").Find(&tournaments)
	return tournaments
}

// GetByHandicapSet returns tournaments using a specific handicap set
func (r *TournamentRepo) GetByHandicapSet(handicapSetID uint) []models.Tournament {
	var tournaments []models.Tournament
	r.db.Where("handicap_set_id = ?", handicapSetID).Order("start_date DESC").Find(&tournaments)
	return tournaments
}