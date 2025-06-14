package tournamentHelper

import (
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
	"sort"
)

func CreateTournamentArchers(db *gorm.DB, tournamentID uint, archerIDs []uint) error {
	var tournament models.Tournament
	print("Creating tournament archers for tournament ID: ", tournamentID, "\n")
	if err := db.Preload("HandicapSet").First(&tournament, tournamentID).Error; err != nil {
		return err
	}

	for _, archerID := range archerIDs {
		var archer models.Archer
		if err := db.Preload("BowClass").First(&archer, archerID).Error; err != nil {
			return err
		}
		if archer.BowClassID == 0 {
			return gorm.ErrRecordNotFound
		}

		score := models.Score{
			ArcherID:     archerID,
			TournamentID: tournamentID,
			TotalScore:   0, // Initialize with 0, will be updated later
			Ranking:      0, // Initialize with 0, will be updated later
			Score:        0, // Initialize with 0, will be updated later
		}
		db.Create(&score)

		var handicapEntry models.HandicapEntry
		if err := db.Where("handicap_set_id = ? AND bow_class_id = ?", tournament.HandicapSetID, archer.BowClassID).First(&handicapEntry).Error; err != nil {
			return err
		}

		tournamentArcher := models.TournamentArcher{
			ArcherID:        archerID,
			Archer:          archer,
			TournamentID:    tournamentID,
			Tournament:      tournament,
			BowClassID:      archer.BowClassID,
			BowClass:        archer.BowClass,
			ScoreID:         score.ID,
			Score:           score,
			HandicapEntryID: handicapEntry.ID,
			HandicapEntry:   handicapEntry,
		}

		if err := db.Create(&tournamentArcher).Error; err != nil {
			return err
		}

	}

	return nil
}

// RecalculateRankings recalculates the rankings of archers in a tournament based on their total scores.
func RecalculateRankings(db *gorm.DB, tournamentID uint) error {
	var scores []models.Score
	if err := db.Where("tournament_id = ?", tournamentID).Order("total_score desc").Find(&scores).Error; err != nil {
		return err
	}

	// Sort scores by TotalScore in descending order
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].TotalScore > scores[j].TotalScore
	})

	// Update rankings
	for i, score := range scores {
		score.Ranking = uint(i + 1)
		if err := db.Save(&score).Error; err != nil {
			return err
		}
	}

	return nil
}
