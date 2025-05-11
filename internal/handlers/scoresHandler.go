package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
		response = append(response, ArcherResponse{
			ID:             archer.ID,
			Name:           archer.FirstName + " " + archer.LastName,
			HandycapFactor: 0.78, // Assuming this is a placeholder
			BowClassName:   archer.BowClass.Name,
			BowClassCode:   archer.BowClass.Code,
		})
	}

	c.JSON(http.StatusOK, response)
}
