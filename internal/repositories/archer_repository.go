package repositories

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/linuxswords/hat/internal/models"
)

type ArcherRepository struct {
	dataFile string
	archers  []models.Archer
	nextID   int
}

// NewArcherRepository creates a new archer repository that implements ArcherRepository
func NewArcherRepository() *ArcherRepository {
	repo := &ArcherRepository{
		dataFile: "doc/data/archers/archers.json",
		archers:  []models.Archer{},
		nextID:   1,
	}
	repo.loadFromFile()
	return repo
}

// loadFromFile loads archers from JSON file
func (r *ArcherRepository) loadFromFile() error {
	data, err := os.ReadFile(r.dataFile)
	if err != nil {
		// If file doesn't exist, start with empty list
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &r.archers)
	if err != nil {
		return err
	}

	// Set nextID to the highest ID + 1
	for _, archer := range r.archers {
		if archer.ID >= r.nextID {
			r.nextID = archer.ID + 1
		}
	}

	return nil
}

// saveToFile saves archers to JSON file
func (r *ArcherRepository) saveToFile() error {
	data, err := json.MarshalIndent(r.archers, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.dataFile, data, 0644)
}

// GetAll returns all archers
func (r *ArcherRepository) GetAll() []models.Archer {
	return r.archers
}

// GetByID returns an archer by ID
func (r *ArcherRepository) GetByID(id int) (*models.Archer, error) {
	for i := range r.archers {
		if r.archers[i].ID == id {
			return &r.archers[i], nil
		}
	}
	return nil, errors.New("archer not found")
}

// Create creates a new archer
func (r *ArcherRepository) Create(archer models.Archer) (*models.Archer, error) {
	archer.ID = r.nextID
	r.nextID++

	r.archers = append(r.archers, archer)

	err := r.saveToFile()
	if err != nil {
		// Rollback the addition if save fails
		r.archers = r.archers[:len(r.archers)-1]
		r.nextID--
		return nil, err
	}

	return &archer, nil
}

// Update updates an existing archer
func (r *ArcherRepository) Update(id int, updatedArcher models.Archer) (*models.Archer, error) {
	for i := range r.archers {
		if r.archers[i].ID == id {
			// Preserve the original ID
			updatedArcher.ID = id
			r.archers[i] = updatedArcher

			err := r.saveToFile()
			if err != nil {
				return nil, err
			}

			return &r.archers[i], nil
		}
	}
	return nil, errors.New("archer not found")
}

// Delete deletes an archer by ID
func (r *ArcherRepository) Delete(id int) error {
	for i := range r.archers {
		if r.archers[i].ID == id {
			// Remove the archer from slice
			r.archers = append(r.archers[:i], r.archers[i+1:]...)

			err := r.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("archer not found")
}

// GetByBowClass returns all archers in a specific bow class
func (r *ArcherRepository) GetByBowClass(bowClass string) []models.Archer {
	var result []models.Archer
	for _, archer := range r.archers {
		if archer.BowClass == bowClass {
			result = append(result, archer)
		}
	}
	return result
}

// GetByGender returns all archers of a specific gender
func (r *ArcherRepository) GetByGender(gender string) []models.Archer {
	var result []models.Archer
	for _, archer := range r.archers {
		if archer.Gender == gender {
			result = append(result, archer)
		}
	}
	return result
}

