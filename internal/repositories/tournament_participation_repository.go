package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// TournamentParticipationRepo implements TournamentParticipationRepository using GORM
type TournamentParticipationRepo struct {
	db *gorm.DB
}

// NewTournamentParticipationRepository creates a new database-backed tournament participation repository
func NewTournamentParticipationRepository(db *gorm.DB) *TournamentParticipationRepo {
	return &TournamentParticipationRepo{db: db}
}

// GetAll returns all tournament participations
func (r *TournamentParticipationRepo) GetAll() []models.TournamentParticipation {
	var participations []models.TournamentParticipation
	r.db.Preload("Tournament").Preload("Archer").Order("created_at DESC").Find(&participations)
	return participations
}

// GetByTournamentID returns all participations for a tournament
func (r *TournamentParticipationRepo) GetByTournamentID(tournamentID uint) []models.TournamentParticipation {
	var participations []models.TournamentParticipation
	r.db.Where("tournament_id = ?", tournamentID).
		Preload("Tournament").
		Preload("Archer").
		Order("created_at ASC").
		Find(&participations)
	return participations
}

// GetByArcherID returns all participations for an archer
func (r *TournamentParticipationRepo) GetByArcherID(archerID uint) []models.TournamentParticipation {
	var participations []models.TournamentParticipation
	r.db.Where("archer_id = ?", archerID).
		Preload("Tournament").
		Preload("Archer").
		Order("created_at DESC").
		Find(&participations)
	return participations
}

// GetByTournamentAndArcher returns a specific participation
func (r *TournamentParticipationRepo) GetByTournamentAndArcher(tournamentID, archerID uint) (*models.TournamentParticipation, error) {
	var participation models.TournamentParticipation
	result := r.db.Where("tournament_id = ? AND archer_id = ?", tournamentID, archerID).
		Preload("Tournament").
		Preload("Archer").
		First(&participation)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("participation for archer %d in tournament %d not found", archerID, tournamentID)
		}
		return nil, result.Error
	}
	return &participation, nil
}

// Create creates a new tournament participation
func (r *TournamentParticipationRepo) Create(participation models.TournamentParticipation) (*models.TournamentParticipation, error) {
	// Check if participation already exists
	var existing models.TournamentParticipation
	result := r.db.Where("tournament_id = ? AND archer_id = ?", participation.TournamentID, participation.ArcherID).First(&existing)
	if result.Error == nil {
		return nil, fmt.Errorf("archer %d is already registered for tournament %d", participation.ArcherID, participation.TournamentID)
	}
	
	result = r.db.Create(&participation)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Reload with associations
	r.db.Preload("Tournament").Preload("Archer").First(&participation, participation.ID)
	return &participation, nil
}

// Update updates an existing tournament participation
func (r *TournamentParticipationRepo) Update(id uint, updatedParticipation models.TournamentParticipation) (*models.TournamentParticipation, error) {
	var participation models.TournamentParticipation
	result := r.db.First(&participation, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("participation with ID %d not found", id)
		}
		return nil, result.Error
	}
	
	// Update fields
	participation.Status = updatedParticipation.Status
	
	result = r.db.Save(&participation)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Reload with associations
	r.db.Preload("Tournament").Preload("Archer").First(&participation, participation.ID)
	return &participation, nil
}

// Delete deletes a tournament participation by ID
func (r *TournamentParticipationRepo) Delete(id uint) error {
	result := r.db.Delete(&models.TournamentParticipation{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("participation with ID %d not found", id)
	}
	return nil
}

// DeleteByTournamentAndArcher deletes a participation by tournament and archer IDs
func (r *TournamentParticipationRepo) DeleteByTournamentAndArcher(tournamentID, archerID uint) error {
	result := r.db.Where("tournament_id = ? AND archer_id = ?", tournamentID, archerID).Delete(&models.TournamentParticipation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("participation for archer %d in tournament %d not found", archerID, tournamentID)
	}
	return nil
}

// GetArcherCountByTournament returns the number of archers in a tournament
func (r *TournamentParticipationRepo) GetArcherCountByTournament(tournamentID uint) int64 {
	var count int64
	r.db.Model(&models.TournamentParticipation{}).Where("tournament_id = ?", tournamentID).Count(&count)
	return count
}

// GetTournamentCountByArcher returns the number of tournaments an archer has participated in
func (r *TournamentParticipationRepo) GetTournamentCountByArcher(archerID uint) int64 {
	var count int64
	r.db.Model(&models.TournamentParticipation{}).Where("archer_id = ?", archerID).Count(&count)
	return count
}