package repositories

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/linuxswords/hat/internal/models"
)

type HandicapRepository struct {
	dataFile     string
	handicapSets []models.HandicapSet
	handicaps    []models.Handicap
	nextID       int
}

// NewHandicapRepository creates a new handicap repository
func NewHandicapRepository() *HandicapRepository {
	repo := &HandicapRepository{
		dataFile:     "doc/data/handicaps/handicapset_h1.json",
		handicapSets: []models.HandicapSet{},
		handicaps:    []models.Handicap{},
		nextID:       1,
	}
	repo.loadFromFile()
	return repo
}

// loadFromFile loads handicaps from JSON file
func (r *HandicapRepository) loadFromFile() error {
	data, err := os.ReadFile(r.dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var jsonData map[string]map[string]float64
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	// Create a default handicap set from the loaded data
	if handicapData, exists := jsonData["Handicap Set"]; exists {
		handicapSet := models.HandicapSet{
			ID:          1,
			Name:        "H1 - Default Set",
			Description: "Default handicap set loaded from handicapset_h1.json",
			IsActive:    true,
		}
		r.handicapSets = []models.HandicapSet{handicapSet}

		// Convert the map to handicap entries
		handicapID := 1
		for bowClass, factor := range handicapData {
			handicap := models.Handicap{
				ID:         handicapID,
				BowClassID: bowClass,
				Factor:     factor,
				SetID:      1,
			}
			r.handicaps = append(r.handicaps, handicap)
			handicapID++
		}
		r.nextID = handicapID
	}

	return nil
}

// GetAllSets returns all handicap sets
func (r *HandicapRepository) GetAllSets() []models.HandicapSet {
	return r.handicapSets
}

// GetSetByID returns a handicap set by ID
func (r *HandicapRepository) GetSetByID(id int) (*models.HandicapSet, error) {
	for i := range r.handicapSets {
		if r.handicapSets[i].ID == id {
			return &r.handicapSets[i], nil
		}
	}
	return nil, errors.New("handicap set not found")
}

// GetHandicapsBySetID returns all handicaps for a specific set
func (r *HandicapRepository) GetHandicapsBySetID(setID int) []models.Handicap {
	var result []models.Handicap
	for _, handicap := range r.handicaps {
		if handicap.SetID == setID {
			result = append(result, handicap)
		}
	}
	return result
}

// GetHandicapByBowClass returns handicap factor for a bow class in a specific set
func (r *HandicapRepository) GetHandicapByBowClass(setID int, bowClassID string) (*models.Handicap, error) {
	for _, handicap := range r.handicaps {
		if handicap.SetID == setID && handicap.BowClassID == bowClassID {
			return &handicap, nil
		}
	}
	return nil, errors.New("handicap not found for bow class")
}

// GetActiveSet returns the currently active handicap set
func (r *HandicapRepository) GetActiveSet() (*models.HandicapSet, error) {
	for i := range r.handicapSets {
		if r.handicapSets[i].IsActive {
			return &r.handicapSets[i], nil
		}
	}
	return nil, errors.New("no active handicap set found")
}