package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// DBArcherRepository implements ArcherRepository using GORM
type DBArcherRepository struct {
	db *gorm.DB
}

// NewDBArcherRepository creates a new database-backed archer repository
func NewDBArcherRepository(db *gorm.DB) *DBArcherRepository {
	return &DBArcherRepository{db: db}
}

// GetAll returns all archers
func (r *DBArcherRepository) GetAll() []models.Archer {
	var archers []models.Archer
	r.db.Find(&archers)
	return archers
}

// GetByID returns an archer by ID
func (r *DBArcherRepository) GetByID(id uint) (*models.Archer, error) {
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
func (r *DBArcherRepository) Create(archer models.Archer) (*models.Archer, error) {
	result := r.db.Create(&archer)
	if result.Error != nil {
		return nil, result.Error
	}
	return &archer, nil
}

// Update updates an existing archer
func (r *DBArcherRepository) Update(id uint, updatedArcher models.Archer) (*models.Archer, error) {
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
func (r *DBArcherRepository) Delete(id uint) error {
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
func (r *DBArcherRepository) GetByBowClass(bowClass string) []models.Archer {
	var archers []models.Archer
	r.db.Where("bow_class = ?", bowClass).Find(&archers)
	return archers
}

// GetByGender returns archers filtered by gender
func (r *DBArcherRepository) GetByGender(gender string) []models.Archer {
	var archers []models.Archer
	r.db.Where("gender = ?", gender).Find(&archers)
	return archers
}