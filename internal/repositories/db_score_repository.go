package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// DBScoreRepository implements ScoreRepository using GORM
type DBScoreRepository struct {
	db *gorm.DB
}

// NewDBScoreRepository creates a new database-backed score repository
func NewDBScoreRepository(db *gorm.DB) *DBScoreRepository {
	return &DBScoreRepository{db: db}
}

// GetAll returns all scores
func (r *DBScoreRepository) GetAll() []models.Score {
	var scores []models.Score
	r.db.Preload("Archer").Preload("Tournament").Order("created_at DESC").Find(&scores)
	return scores
}

// GetByTournamentID returns all scores for a tournament
func (r *DBScoreRepository) GetByTournamentID(tournamentID uint) []models.Score {
	var scores []models.Score
	r.db.Where("tournament_id = ?", tournamentID).
		Preload("Archer").
		Preload("Tournament").
		Order("raw_score DESC").
		Find(&scores)
	return scores
}

// GetByTournamentIDSorted returns scores for a tournament sorted by adjusted score (or raw score if no handicap)
func (r *DBScoreRepository) GetByTournamentIDSorted(tournamentID uint) []models.Score {
	var scores []models.Score
	r.db.Where("tournament_id = ?", tournamentID).
		Preload("Archer").
		Preload("Tournament").
		Order("COALESCE(adjusted_score, raw_score) DESC").
		Find(&scores)
	return scores
}

// GetByArcherID returns all scores for an archer
func (r *DBScoreRepository) GetByArcherID(archerID uint) []models.Score {
	var scores []models.Score
	r.db.Where("archer_id = ?", archerID).
		Preload("Archer").
		Preload("Tournament").
		Order("created_at DESC").
		Find(&scores)
	return scores
}

// GetByTournamentAndArcher returns a score for a specific archer in a tournament
func (r *DBScoreRepository) GetByTournamentAndArcher(tournamentID, archerID uint) (*models.Score, error) {
	var score models.Score
	result := r.db.Where("tournament_id = ? AND archer_id = ?", tournamentID, archerID).
		Preload("Archer").
		Preload("Tournament").
		First(&score)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("score for archer %d in tournament %d not found", archerID, tournamentID)
		}
		return nil, result.Error
	}
	return &score, nil
}

// Create creates a new score
func (r *DBScoreRepository) Create(score models.Score) (*models.Score, error) {
	result := r.db.Create(&score)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Reload with associations
	r.db.Preload("Archer").Preload("Tournament").First(&score, score.ID)
	return &score, nil
}

// Update updates an existing score
func (r *DBScoreRepository) Update(id uint, updatedScore models.Score) (*models.Score, error) {
	var score models.Score
	result := r.db.First(&score, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("score with ID %d not found", id)
		}
		return nil, result.Error
	}
	
	// Update fields
	score.RawScore = updatedScore.RawScore
	score.AdjustedScore = updatedScore.AdjustedScore
	score.BowClassID = updatedScore.BowClassID
	score.EnteredBy = updatedScore.EnteredBy
	
	result = r.db.Save(&score)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Reload with associations
	r.db.Preload("Archer").Preload("Tournament").First(&score, score.ID)
	return &score, nil
}

// UpdateByTournamentAndArcher updates a score by tournament and archer IDs
func (r *DBScoreRepository) UpdateByTournamentAndArcher(tournamentID, archerID uint, updatedScore models.Score) (*models.Score, error) {
	var score models.Score
	result := r.db.Where("tournament_id = ? AND archer_id = ?", tournamentID, archerID).First(&score)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("score for archer %d in tournament %d not found", archerID, tournamentID)
		}
		return nil, result.Error
	}
	
	// Update fields
	score.RawScore = updatedScore.RawScore
	score.AdjustedScore = updatedScore.AdjustedScore
	score.BowClassID = updatedScore.BowClassID
	score.EnteredBy = updatedScore.EnteredBy
	
	result = r.db.Save(&score)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Reload with associations
	r.db.Preload("Archer").Preload("Tournament").First(&score, score.ID)
	return &score, nil
}

// Delete deletes a score by ID
func (r *DBScoreRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Score{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("score with ID %d not found", id)
	}
	return nil
}

// DeleteByTournamentAndArcher deletes a score by tournament and archer IDs
func (r *DBScoreRepository) DeleteByTournamentAndArcher(tournamentID, archerID uint) error {
	result := r.db.Where("tournament_id = ? AND archer_id = ?", tournamentID, archerID).Delete(&models.Score{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("score for archer %d in tournament %d not found", archerID, tournamentID)
	}
	return nil
}

// GetScoreCountByTournament returns the number of scores for a tournament
func (r *DBScoreRepository) GetScoreCountByTournament(tournamentID uint) int64 {
	var count int64
	r.db.Model(&models.Score{}).Where("tournament_id = ?", tournamentID).Count(&count)
	return count
}

// CalculateAdjustedScore calculates adjusted score using handicap factor
func (r *DBScoreRepository) CalculateAdjustedScore(rawScore int, handicapFactor float64) float64 {
	return float64(rawScore) * handicapFactor
}