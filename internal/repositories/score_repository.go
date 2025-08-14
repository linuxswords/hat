package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"

	"github.com/linuxswords/hat/internal/models"
)

type ScoreRepository struct {
	dataFile string
	scores   []models.Score
	nextID   int
}

// NewScoreRepository creates a new score repository
func NewScoreRepository() *ScoreRepository {
	repo := &ScoreRepository{
		dataFile: "doc/data/tournaments/scores.json",
		scores:   []models.Score{},
		nextID:   1,
	}
	repo.loadFromFile()
	return repo
}

// loadFromFile loads scores from JSON file
func (r *ScoreRepository) loadFromFile() error {
	data, err := os.ReadFile(r.dataFile)
	if err != nil {
		// If file doesn't exist, start with empty list
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &r.scores)
	if err != nil {
		return err
	}

	// Set nextID to the highest ID + 1
	for _, score := range r.scores {
		if score.ID >= r.nextID {
			r.nextID = score.ID + 1
		}
	}

	return nil
}

// saveToFile saves scores to JSON file
func (r *ScoreRepository) saveToFile() error {
	// Ensure directory exists
	err := os.MkdirAll("doc/data/tournaments", 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.scores, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.dataFile, data, 0644)
}

// GetAll returns all scores
func (r *ScoreRepository) GetAll() []models.Score {
	return r.scores
}

// GetByTournamentID returns all scores for a specific tournament
func (r *ScoreRepository) GetByTournamentID(tournamentID int) []models.Score {
	var result []models.Score
	for _, score := range r.scores {
		if score.TournamentID == tournamentID {
			result = append(result, score)
		}
	}
	return result
}

// GetByTournamentIDSorted returns all scores for a tournament sorted by adjusted score (or raw if no handicap)
func (r *ScoreRepository) GetByTournamentIDSorted(tournamentID int) []models.Score {
	scores := r.GetByTournamentID(tournamentID)
	
	sort.Slice(scores, func(i, j int) bool {
		// If both have adjusted scores, compare those
		if scores[i].AdjustedScore != nil && scores[j].AdjustedScore != nil {
			return *scores[i].AdjustedScore > *scores[j].AdjustedScore
		}
		// If neither has adjusted scores, compare raw scores
		if scores[i].AdjustedScore == nil && scores[j].AdjustedScore == nil {
			return scores[i].RawScore > scores[j].RawScore
		}
		// If only one has adjusted score, prioritize that one
		return scores[i].AdjustedScore != nil
	})
	
	return scores
}

// GetByArcherID returns all scores for a specific archer
func (r *ScoreRepository) GetByArcherID(archerID int) []models.Score {
	var result []models.Score
	for _, score := range r.scores {
		if score.ArcherID == archerID {
			result = append(result, score)
		}
	}
	return result
}

// GetByTournamentAndArcher returns score for a specific tournament and archer
func (r *ScoreRepository) GetByTournamentAndArcher(tournamentID int, archerID int) (*models.Score, error) {
	for i := range r.scores {
		if r.scores[i].TournamentID == tournamentID && r.scores[i].ArcherID == archerID {
			return &r.scores[i], nil
		}
	}
	return nil, errors.New("score not found")
}

// Create creates a new score
func (r *ScoreRepository) Create(score models.Score) (*models.Score, error) {
	// Check if score already exists for this tournament and archer
	existing, _ := r.GetByTournamentAndArcher(score.TournamentID, score.ArcherID)
	if existing != nil {
		return nil, errors.New("score already exists for this archer in this tournament")
	}

	score.ID = r.nextID
	score.EnteredAt = time.Now()
	r.nextID++

	r.scores = append(r.scores, score)

	err := r.saveToFile()
	if err != nil {
		// Rollback the addition if save fails
		r.scores = r.scores[:len(r.scores)-1]
		r.nextID--
		return nil, err
	}

	return &score, nil
}

// Update updates an existing score
func (r *ScoreRepository) Update(id int, updatedScore models.Score) (*models.Score, error) {
	for i := range r.scores {
		if r.scores[i].ID == id {
			// Preserve the original ID and entry time
			updatedScore.ID = id
			updatedScore.EnteredAt = r.scores[i].EnteredAt
			r.scores[i] = updatedScore

			err := r.saveToFile()
			if err != nil {
				return nil, err
			}

			return &r.scores[i], nil
		}
	}
	return nil, errors.New("score not found")
}

// UpdateByTournamentAndArcher updates score by tournament and archer ID
func (r *ScoreRepository) UpdateByTournamentAndArcher(tournamentID int, archerID int, updatedScore models.Score) (*models.Score, error) {
	for i := range r.scores {
		if r.scores[i].TournamentID == tournamentID && r.scores[i].ArcherID == archerID {
			// Preserve the original ID and entry time
			updatedScore.ID = r.scores[i].ID
			updatedScore.EnteredAt = r.scores[i].EnteredAt
			updatedScore.TournamentID = tournamentID
			updatedScore.ArcherID = archerID
			r.scores[i] = updatedScore

			err := r.saveToFile()
			if err != nil {
				return nil, err
			}

			return &r.scores[i], nil
		}
	}
	return nil, errors.New("score not found")
}

// Delete deletes a score by ID
func (r *ScoreRepository) Delete(id int) error {
	for i := range r.scores {
		if r.scores[i].ID == id {
			// Remove the score from slice
			r.scores = append(r.scores[:i], r.scores[i+1:]...)

			err := r.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("score not found")
}

// DeleteByTournamentAndArcher deletes score by tournament and archer ID
func (r *ScoreRepository) DeleteByTournamentAndArcher(tournamentID int, archerID int) error {
	for i := range r.scores {
		if r.scores[i].TournamentID == tournamentID && r.scores[i].ArcherID == archerID {
			// Remove the score from slice
			r.scores = append(r.scores[:i], r.scores[i+1:]...)

			err := r.saveToFile()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("score not found")
}

// GetScoreCountByTournament returns the number of scores entered for a tournament
func (r *ScoreRepository) GetScoreCountByTournament(tournamentID int) int {
	count := 0
	for _, score := range r.scores {
		if score.TournamentID == tournamentID {
			count++
		}
	}
	return count
}

// CalculateAdjustedScore calculates adjusted score using handicap factor
func (r *ScoreRepository) CalculateAdjustedScore(rawScore int, handicapFactor float64) float64 {
	return float64(rawScore) * handicapFactor
}