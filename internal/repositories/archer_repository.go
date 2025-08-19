package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// ArcherRepo implements ArcherRepository using GORM
type ArcherRepo struct {
	db *gorm.DB
}

// NewArcherRepository creates a new database-backed archer repository
func NewArcherRepository(db *gorm.DB) *ArcherRepo {
	return &ArcherRepo{db: db}
}

// GetAll returns all archers
func (r *ArcherRepo) GetAll() []models.Archer {
	var archers []models.Archer
	r.db.Find(&archers)
	return archers
}

// GetByID returns an archer by ID
func (r *ArcherRepo) GetByID(id uint) (*models.Archer, error) {
	var archer models.Archer
	result := r.db.First(&archer, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("archer with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &archer, nil
}

// Create creates a new archer
func (r *ArcherRepo) Create(archer models.Archer) (*models.Archer, error) {
	result := r.db.Create(&archer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &archer, nil
}

// Update updates an existing archer
func (r *ArcherRepo) Update(id uint, updatedArcher models.Archer) (*models.Archer, error) {
	var archer models.Archer
	result := r.db.First(&archer, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("archer with ID %d not found", id)
		}
		return nil, result.Error
	}
	
	// Update fields
	archer.Name = updatedArcher.Name
	archer.Gender = updatedArcher.Gender
	archer.BowClass = updatedArcher.BowClass
	archer.Email = updatedArcher.Email
	
	result = r.db.Save(&archer)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &archer, nil
}

// Delete deletes an archer by ID
func (r *ArcherRepo) Delete(id uint) error {
	result := r.db.Delete(&models.Archer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("archer with ID %d not found", id)
	}
	return nil
}

// GetByBowClass returns archers filtered by bow class
func (r *ArcherRepo) GetByBowClass(bowClass string) []models.Archer {
	var archers []models.Archer
	r.db.Where("bow_class = ?", bowClass).Find(&archers)
	return archers
}

// GetByGender returns archers filtered by gender
func (r *ArcherRepo) GetByGender(gender string) []models.Archer {
	var archers []models.Archer
	r.db.Where("gender = ?", gender).Find(&archers)
	return archers
}