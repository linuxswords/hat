package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// ScoreRepo implements ScoreRepository using GORM
type ScoreRepo struct {
	db *gorm.DB
}

// NewScoreRepository creates a new database-backed score repository
func NewScoreRepository(db *gorm.DB) *ScoreRepo {
	return &ScoreRepo{db: db}
}

// GetAll returns all scores
func (r *ScoreRepo) GetAll() []models.Score {
	var scores []models.Score
	r.db.Preload("Archer").Preload("Tournament").Order("created_at DESC").Find(&scores)
	return scores
}

// GetByTournamentID returns all scores for a tournament
func (r *ScoreRepo) GetByTournamentID(tournamentID uint) []models.Score {
	var scores []models.Score
	r.db.Where("tournament_id = ?", tournamentID).
		Preload("Archer").
		Preload("Tournament").
		Order("raw_score DESC").
		Find(&scores)
	return scores
}

// GetByTournamentIDSorted returns scores for a tournament sorted by adjusted score (or raw score if no handicap)
func (r *ScoreRepo) GetByTournamentIDSorted(tournamentID uint) []models.Score {
	var scores []models.Score
	r.db.Where("tournament_id = ?", tournamentID).
		Preload("Archer").
		Preload("Tournament").
		Order("COALESCE(adjusted_score, raw_score) DESC").
		Find(&scores)
	return scores
}

// GetByArcherID returns all scores for an archer
func (r *ScoreRepo) GetByArcherID(archerID uint) []models.Score {
	var scores []models.Score
	r.db.Where("archer_id = ?", archerID).
		Preload("Archer").
		Preload("Tournament").
		Order("created_at DESC").
		Find(&scores)
	return scores
}

// GetByTournamentAndArcher returns a score for a specific archer in a tournament
func (r *ScoreRepo) GetByTournamentAndArcher(tournamentID, archerID uint) (*models.Score, error) {
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
func (r *ScoreRepo) Create(score models.Score) (*models.Score, error) {
	result := r.db.Create(&score)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Reload with associations
	r.db.Preload("Archer").Preload("Tournament").First(&score, score.ID)
	return &score, nil
}

// Update updates an existing score
func (r *ScoreRepo) Update(id uint, updatedScore models.Score) (*models.Score, error) {
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
func (r *ScoreRepo) UpdateByTournamentAndArcher(tournamentID, archerID uint, updatedScore models.Score) (*models.Score, error) {
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
func (r *ScoreRepo) Delete(id uint) error {
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
func (r *ScoreRepo) DeleteByTournamentAndArcher(tournamentID, archerID uint) error {
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
func (r *ScoreRepo) GetScoreCountByTournament(tournamentID uint) int64 {
	var count int64
	r.db.Model(&models.Score{}).Where("tournament_id = ?", tournamentID).Count(&count)
	return count
}

// CalculateAdjustedScore calculates adjusted score using handicap factor
func (r *ScoreRepo) CalculateAdjustedScore(rawScore int, handicapFactor float64) float64 {
	return float64(rawScore) * handicapFactor
}