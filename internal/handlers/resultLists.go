package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/linuxswords/hat/internal/pdf"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ShowResultLists(c *gin.Context) {
	// Fetch tournaments from the database
	var scores []models.Score
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&scores)
	tournamentIdMap := make(map[uint]bool)
	var tournamentIds []uint
	for _, s := range scores {
		if _, exists := tournamentIdMap[s.TournamentID]; !exists {
			tournamentIdMap[s.TournamentID] = true
			tournamentIds = append(tournamentIds, s.TournamentID)
		}
	}
	var tournaments []models.Tournament
	db.Where("ID IN ?", tournamentIds).Find(&tournaments)

	c.HTML(http.StatusOK, "resultLists.tmpl", gin.H{
		"Title":       "HAT - Result Lists",
		"Content":     "Tournament Result Lists",
		"Tournaments": tournaments,
	})
}

func DownloadTournamentPDF(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var scores []models.Score
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("tournament_id = ?", id).Find(&scores).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scores not found!"})
		return
	}
	pdfData, err := pdf.CreatePDF(db, scores)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	// Set the response headers for PDF download
	c.Header("Content-Disposition", "attachment; filename=tournament_results.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
}
