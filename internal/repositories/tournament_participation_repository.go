package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/linuxswords/hat/internal/models"
)

type TournamentParticipationRepository struct {
	dataFile       string
	participations []models.TournamentParticipation
	nextID         int
}

// NewTournamentParticipationRepository creates a new tournament participation repository
func NewTournamentParticipationRepository() *TournamentParticipationRepository {
	repo := &TournamentParticipationRepository{
		dataFile:       "doc/data/tournaments/participations.json",
		participations: []models.TournamentParticipation{},
		nextID:         1,
	}
	repo.loadFromFile()
	return repo
}

// loadFromFile loads participations from JSON file
func (r *TournamentParticipationRepository) loadFromFile() error {
	data, err := os.ReadFile(r.dataFile)
	if err != nil {
		// If file doesn't exist, start with empty list
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &r.participations)
	if err != nil {
		return err
	}

	// Set nextID to the highest ID + 1
	for _, participation := range r.participations {
		if participation.ID >= r.nextID {
			r.nextID = participation.ID + 1
		}
	}

	return nil
}

// saveToFile saves participations to JSON file
func (r *TournamentParticipationRepository) saveToFile() error {
	// Ensure directory exists
	err := os.MkdirAll("doc/data/tournaments", 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.participations, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.dataFile, data, 0644)
}

// GetAll returns all participations
func (r *TournamentParticipationRepository) GetAll() []models.TournamentParticipation {
	return r.participations
}

// GetByTournamentID returns all participations for a specific tournament
func (r *TournamentParticipationRepository) GetByTournamentID(tournamentID int) []models.TournamentParticipation {
	var result []models.TournamentParticipation
	for _, participation := range r.participations {
		if participation.TournamentID == tournamentID {
			result = append(result, participation)
		}
	}
	return result
}

// GetByArcherID returns all participations for a specific archer
func (r *TournamentParticipationRepository) GetByArcherID(archerID int) []models.TournamentParticipation {
	var result []models.TournamentParticipation
	for _, participation := range r.participations {
		if participation.ArcherID == archerID {
			result = append(result, participation)
		}
	}
	return result
}

// GetByTournamentAndArcher returns participation for a specific tournament and archer
func (r *TournamentParticipationRepository) GetByTournamentAndArcher(tournamentID int, archerID int) (*models.TournamentParticipation, error) {
	for i := range r.participations {
		if r.participations[i].TournamentID == tournamentID && r.participations[i].ArcherID == archerID {
			return &r.participations[i], nil
		}
	}
	return nil, errors.New("participation not found")
}

// Create creates a new participation
func (r *TournamentParticipationRepository) Create(participation models.TournamentParticipation) (*models.TournamentParticipation, error) {
	// Check if archer is already registered for this tournament
	existing, _ := r.GetByTournamentAndArcher(participation.TournamentID, participation.ArcherID)
	if existing != nil {
		return nil, errors.New("archer is already registered for this tournament")
	}

	participation.ID = r.nextID
	participation.RegisteredAt = time.Now()
	if participation.Status == "" {
		participation.Status = "registered"
	}
	r.nextID++

	r.participations = append(r.participations, participation)

	err := r.saveToFile()
	if err != nil {
		// Rollback the addition if save fails
		r.participations = r.participations[:len(r.participations)-1]
		r.nextID--
		return nil, err
	}

	return &participation, nil
}

// Update updates an existing participation
func (r *TournamentParticipationRepository) Update(id int, updatedParticipation models.TournamentParticipation) (*models.TournamentParticipation, error) {
	for i := range r.participations {
		if r.participations[i].ID == id {
			// Preserve the original ID and registration time
			updatedParticipation.ID = id
			updatedParticipation.RegisteredAt = r.participations[i].RegisteredAt
			r.participations[i] = updatedParticipation

			err := r.saveToFile()
			if err != nil {
				return nil, err
			}

			return &r.participations[i], nil
		}
	}
	return nil, errors.New("participation not found")
}

// Delete deletes a participation by ID
func (r *TournamentParticipationRepository) Delete(id int) error {
	for i := range r.participations {
		if r.participations[i].ID == id {
			// Remove the participation from slice
			r.participations = append(r.participations[:i], r.participations[i+1:]...)

			err := r.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("participation not found")
}

// DeleteByTournamentAndArcher deletes participation by tournament and archer ID
func (r *TournamentParticipationRepository) DeleteByTournamentAndArcher(tournamentID int, archerID int) error {
	for i := range r.participations {
		if r.participations[i].TournamentID == tournamentID && r.participations[i].ArcherID == archerID {
			// Remove the participation from slice
			r.participations = append(r.participations[:i], r.participations[i+1:]...)

			err := r.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("participation not found")
}

// GetArcherCountByTournament returns the number of archers registered for a tournament
func (r *TournamentParticipationRepository) GetArcherCountByTournament(tournamentID int) int {
	count := 0
	for _, participation := range r.participations {
		if participation.TournamentID == tournamentID {
			count++
		}
	}
	return count
}

// GetTournamentCountByArcher returns the number of tournaments an archer is registered for
func (r *TournamentParticipationRepository) GetTournamentCountByArcher(archerID int) int {
	count := 0
	for _, participation := range r.participations {
		if participation.ArcherID == archerID {
			count++
		}
	}
	return count
}