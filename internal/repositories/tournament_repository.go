package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/linuxswords/hat/internal/models"
)

type TournamentRepository struct {
	dataFile    string
	tournaments []models.Tournament
	nextID      int
}

// NewTournamentRepository creates a new tournament repository
func NewTournamentRepository() *TournamentRepository {
	repo := &TournamentRepository{
		dataFile:    "doc/data/tournaments/tournaments.json",
		tournaments: []models.Tournament{},
		nextID:      1,
	}
	repo.loadFromFile()
	return repo
}

// loadFromFile loads tournaments from JSON file
func (r *TournamentRepository) loadFromFile() error {
	data, err := os.ReadFile(r.dataFile)
	if err != nil {
		// If file doesn't exist, start with empty list
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &r.tournaments)
	if err != nil {
		return err
	}

	// Set nextID to the highest ID + 1
	for _, tournament := range r.tournaments {
		if tournament.ID >= r.nextID {
			r.nextID = tournament.ID + 1
		}
	}

	return nil
}

// saveToFile saves tournaments to JSON file
func (r *TournamentRepository) saveToFile() error {
	// Ensure directory exists
	err := os.MkdirAll("doc/data/tournaments", 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.tournaments, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.dataFile, data, 0644)
}

// GetAll returns all tournaments
func (r *TournamentRepository) GetAll() []models.Tournament {
	return r.tournaments
}

// GetByID returns a tournament by ID
func (r *TournamentRepository) GetByID(id int) (*models.Tournament, error) {
	for i := range r.tournaments {
		if r.tournaments[i].ID == id {
			return &r.tournaments[i], nil
		}
	}
	return nil, errors.New("tournament not found")
}

// Create creates a new tournament
func (r *TournamentRepository) Create(tournament models.Tournament) (*models.Tournament, error) {
	tournament.ID = r.nextID
	r.nextID++

	r.tournaments = append(r.tournaments, tournament)

	err := r.saveToFile()
	if err != nil {
		// Rollback the addition if save fails
		r.tournaments = r.tournaments[:len(r.tournaments)-1]
		r.nextID--
		return nil, err
	}

	return &tournament, nil
}

// Update updates an existing tournament
func (r *TournamentRepository) Update(id int, updatedTournament models.Tournament) (*models.Tournament, error) {
	for i := range r.tournaments {
		if r.tournaments[i].ID == id {
			// Preserve the original ID
			updatedTournament.ID = id
			r.tournaments[i] = updatedTournament

			err := r.saveToFile()
			if err != nil {
				return nil, err
			}

			return &r.tournaments[i], nil
		}
	}
	return nil, errors.New("tournament not found")
}

// Delete deletes a tournament by ID
func (r *TournamentRepository) Delete(id int) error {
	for i := range r.tournaments {
		if r.tournaments[i].ID == id {
			// Remove the tournament from slice
			r.tournaments = append(r.tournaments[:i], r.tournaments[i+1:]...)

			err := r.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("tournament not found")
}

// GetUpcoming returns tournaments that haven't ended yet
func (r *TournamentRepository) GetUpcoming() []models.Tournament {
	var upcoming []models.Tournament
	now := time.Now()
	
	for _, tournament := range r.tournaments {
		if tournament.EndDate.After(now) {
			upcoming = append(upcoming, tournament)
		}
	}
	
	return upcoming
}

// GetCurrent returns tournaments that are currently running
func (r *TournamentRepository) GetCurrent() []models.Tournament {
	var current []models.Tournament
	now := time.Now()
	
	for _, tournament := range r.tournaments {
		if tournament.StartDate.Before(now) && tournament.EndDate.After(now) {
			current = append(current, tournament)
		}
	}
	
	return current
}

// GetByHandicapSet returns tournaments using a specific handicap set
func (r *TournamentRepository) GetByHandicapSet(handicapSetID int) []models.Tournament {
	var result []models.Tournament
	
	for _, tournament := range r.tournaments {
		if tournament.HandicapSetID != nil && *tournament.HandicapSetID == handicapSetID {
			result = append(result, tournament)
		}
	}
	
	return result
}