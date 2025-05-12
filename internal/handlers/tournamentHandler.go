package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func ShowTournamentsPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tournaments []models.Tournament
	db.Preload("HandycapSet").Preload("Archers").Find(&tournaments)
	var handycapSets []models.HandycapSet
	db.Find(&handycapSets)
	var archers []models.Archer
	db.Find(&archers)
	c.HTML(http.StatusOK, "tournaments.tmpl", gin.H{
		"Title":        "Tournaments",
		"Tournaments":  tournaments,
		"HandycapSets": handycapSets,
		"Archers":      archers,
	})
}

func AddTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tournament models.Tournament
	if err := c.ShouldBind(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tournament.TournamentType = c.PostForm("TournamentType")
	handycapSetID, _ := strconv.Atoi(c.PostForm("handycapSetID"))
	var handycapSet models.HandycapSet
	if err := db.First(&handycapSet, handycapSetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HandycapSet not found"})
		return
	}

	tournament.HandycapSetID = handycapSet.ID
	tournament.HandycapSet = handycapSet

	tournament.Date, _ = time.Parse("2006-01-02", c.PostForm("date"))
	archerIDs := c.PostFormArray("archers")
	for _, archerID := range archerIDs {
		id, _ := strconv.Atoi(archerID)
		var archer models.Archer
		if err := db.First(&archer, id).Error; err == nil {
			db.Model(&tournament).Association("Archers").Append(&archer)
		}
	}
	db.Create(&tournament)
	c.Redirect(http.StatusSeeOther, "/tournaments")
}

func UpdateTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var tournament models.Tournament
	if err := db.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tournament.Date, _ = time.Parse("2006-01-02", c.PostForm("date"))
	db.Save(&tournament)
	c.Status(http.StatusOK)
}

func DeleteTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.Tournament{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/tournaments")
}
