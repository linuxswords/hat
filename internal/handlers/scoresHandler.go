package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/helper"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

// Prepare response with handicap factors
type ArcherResponse struct {
	ID             uint
	Ranking        uint
	Name           string
	BowClassName   string
	BowClassCode   string
	HandicapFactor float64
	TotalScore     float64
	Score          float64
	ScoreID        uint
	TournamentID   uint
}

func createArchersResponse(tournamentArchers []models.TournamentArcher, tournamentId uint) []ArcherResponse {
	var archersData []ArcherResponse
	for _, tournamenArcher := range tournamentArchers {
		factor := 1.0
		if tournamenArcher.HandicapEntry.Value != 0 {
			factor = tournamenArcher.HandicapEntry.Value
		}

		archerResponse := ArcherResponse{
			ID:             tournamenArcher.ID,
			Name:           tournamenArcher.Name(),
			HandicapFactor: factor,
			BowClassName:   tournamenArcher.BowClass.Name,
			BowClassCode:   tournamenArcher.BowClass.Code,
			Ranking:        tournamenArcher.Score.Ranking,
			Score:          tournamenArcher.Score.Score,
			TotalScore:     tournamenArcher.Score.TotalScore,
			TournamentID:   tournamentId,
		}
		archersData = append(archersData, archerResponse)
	}
	// Sort archersData by Ranking
	sort.Slice(archersData, func(i, j int) bool {
		return archersData[i].Ranking < archersData[j].Ranking
	})
	return archersData
}

func ShowTournamentScores(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tournament models.Tournament
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Preload("HandicapSet").First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	var tournamentArchers []models.TournamentArcher
	if err := db.Where("tournament_id = ?", tournament.ID).Preload("HandicapEntry").Preload("Score").Preload("Archer").Preload("BowClass").Find(&tournamentArchers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournament archers"})
		return
	}

	archersData := createArchersResponse(tournamentArchers, tournament.ID)

	c.HTML(http.StatusOK, "scores", gin.H{
		"Title":      "HAT - Tournament Scores",
		"Content":    "Scores for Tournament: " + tournament.Name,
		"Tournament": tournament,
		"Archers":    archersData,
	})
}

func UpdateArcherScore(c *gin.Context) {
	tournamentID, _ := strconv.Atoi(c.Param("tournamentID"))
	tournamentArcherID, _ := strconv.Atoi(c.Param("archerID"))
	db := c.MustGet("db").(*gorm.DB)
	var tournamentArcher models.TournamentArcher
	if err := db.Preload("Score").Preload("Archer").Preload("BowClass").Preload("HandicapEntry").First(&tournamentArcher, tournamentArcherID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament-Archer not found"})
		return
	}
	newScore, err := strconv.ParseFloat(c.PostForm("archerScore"), 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score value"})
	}

	tournamentArcher.Score.Score = newScore
	tournamentArcher.Score.TotalScore = newScore * tournamentArcher.HandicapEntry.Value
	db.Save(&tournamentArcher.Score)

	// recalculate the overall rankings for the tournament
	tournamentHelper.RecalculateRankings(db, tournamentID)

	archerResponse := createArchersResponse([]models.TournamentArcher{tournamentArcher}, uint(tournamentID))

	var tournamentArchers []models.TournamentArcher
	if err := db.Where("tournament_id = ?", tournamentID).Preload("HandicapEntry").Preload("Score").Preload("Archer").Preload("BowClass").Find(&tournamentArchers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournament archers"})
		return
	}
	archersData := createArchersResponse(tournamentArchers, uint(tournamentID))

	c.HTML(http.StatusOK, "archer", archerResponse)
	c.HTML(http.StatusOK, "scoredArchers-oob", gin.H{"Archers": archersData})
}
