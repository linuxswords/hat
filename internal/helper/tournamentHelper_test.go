package tournamentHelper

import (
	"testing"

	"github.com/linuxswords/hat/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRecalculateRankings(t *testing.T) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Score{})

	// Insert test data
	scores := []models.Score{
		{TournamentID: 1, ArcherID: 1, TotalScore: 300, Ranking: 1},
		{TournamentID: 1, ArcherID: 2, TotalScore: 250, Ranking: 1},
		{TournamentID: 1, ArcherID: 3, TotalScore: 275, Ranking: 1},
	}
	for _, score := range scores {
		db.Create(&score)
	}

	// Call the function to test
	err = RecalculateRankings(db, 1)
	assert.NoError(t, err)

	// Verify the results
	var updatedScores []models.Score
	db.Order("ranking asc").Find(&updatedScores)

	assert.Equal(t, uint(1), updatedScores[0].Ranking)
	assert.Equal(t, uint(1), updatedScores[0].ArcherID)

	assert.Equal(t, uint(2), updatedScores[1].Ranking)
	assert.Equal(t, uint(3), updatedScores[1].ArcherID)

	assert.Equal(t, uint(3), updatedScores[2].Ranking)
	assert.Equal(t, uint(2), updatedScores[2].ArcherID)
}
