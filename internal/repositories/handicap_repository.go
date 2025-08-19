package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// HandicapRepo implements HandicapRepository using GORM
type HandicapRepo struct {
	db *gorm.DB
}

// NewHandicapRepository creates a new database-backed handicap repository
func NewHandicapRepository(db *gorm.DB) *HandicapRepo {
	return &HandicapRepo{db: db}
}

// GetAllSets returns all handicap sets
func (r *HandicapRepo) GetAllSets() []models.HandicapSet {
	var sets []models.HandicapSet
	r.db.Order("created_at DESC").Find(&sets)
	return sets
}

// GetSetByID returns a handicap set by ID
func (r *HandicapRepo) GetSetByID(id uint) (*models.HandicapSet, error) {
	var set models.HandicapSet
	result := r.db.First(&set, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("handicap set with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &set, nil
}

// GetHandicapsBySetID returns all handicaps for a specific set
func (r *HandicapRepo) GetHandicapsBySetID(setID uint) []models.Handicap {
	var handicaps []models.Handicap
	r.db.Where("set_id = ?", setID).Order("bow_class_id ASC").Find(&handicaps)
	return handicaps
}

// GetHandicapByBowClass returns a specific handicap for a bow class in a set
func (r *HandicapRepo) GetHandicapByBowClass(setID uint, bowClassID string) (*models.Handicap, error) {
	var handicap models.Handicap
	result := r.db.Where("set_id = ? AND bow_class_id = ?", setID, bowClassID).First(&handicap)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("handicap for bow class %s in set %d not found", bowClassID, setID)
		}
		return nil, result.Error
	}
	return &handicap, nil
}

// GetActiveSet returns the currently active handicap set
func (r *HandicapRepo) GetActiveSet() (*models.HandicapSet, error) {
	var set models.HandicapSet
	result := r.db.Where("is_active = ?", true).First(&set)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no active handicap set found")
		}
		return nil, result.Error
	}
	return &set, nil
}

// CalculateAdjustedScore calculates adjusted score using handicap factor
func (r *HandicapRepo) CalculateAdjustedScore(rawScore int, handicapFactor float64) float64 {
	return float64(rawScore) * handicapFactor
}