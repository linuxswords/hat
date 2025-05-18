package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func ShowScoresPage(c *gin.Context) {
	// Fetch tournaments from the database
	var tournaments []models.Tournament
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&tournaments)

	c.HTML(http.StatusOK, "scores.tmpl", gin.H{
		"Title":       "HAT - Scores",
		"Content":     "Scores Page",
		"Tournaments": tournaments,
	})
}

func SaveScores(c *gin.Context) {
	var request struct {
		TournamentID uint `json:"tournamentId"`
		Scores       []struct {
			ArcherID   uint    `json:"archerId"`
			Score      float64 `json:"score"`
			Ranking    uint    `json:"ranking"`
			TotalScore float64 `json:"totalScore"`
		} `json:"scores"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	for _, scoreData := range request.Scores {
		score := models.Score{
			TournamentID: request.TournamentID,
			ArcherID:     scoreData.ArcherID,
			Score:        scoreData.Score,
			TotalScore:   scoreData.TotalScore,
			Ranking:      scoreData.Ranking,
		}
		db.Create(&score)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scores saved successfully"})
}

func GetHandycapSet(c *gin.Context) {
}

func GetArchers(c *gin.Context) {
	// Fetch archers and their handycap factors from the database
	db := c.MustGet("db").(*gorm.DB)
	var tournament models.Tournament
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	archers := tournament.Archers
	db.Preload("BowClass").Find(&archers)
	hcSet := tournament.HandycapSet
	db.Preload("HandycapEntries").Find(&hcSet)

	fmt.Println("hcEntries", hcSet)

	// Prepare response with handycap factors
	type ArcherResponse struct {
		ID             uint    `json:"id"`
		Name           string  `json:"name"`
		BowClassName   string  `json:"bowclassName"`
		BowClassCode   string  `json:"bowclassCode"`
		HandycapFactor float64 `json:"handycapFactor"`
	}

	var response []ArcherResponse
	for _, archer := range archers {
		factor := 1.0
		hcEntry := hcSet.GetHandycapEntryByBowClass(archer.BowClassID)
		if hcEntry != nil {
			factor = hcEntry.Value
		}
		response = append(response, ArcherResponse{
			ID:             archer.ID,
			Name:           archer.FirstName + " " + archer.LastName,
			HandycapFactor: factor,
			BowClassName:   archer.BowClass.Name,
			BowClassCode:   archer.BowClass.Code,
		})
	}

	c.JSON(http.StatusOK, response)
}
