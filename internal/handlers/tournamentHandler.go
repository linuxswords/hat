package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	tournamentHelper "github.com/linuxswords/hat/internal/helper"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func ShowTournamentsPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tournaments []models.Tournament
	db.Preload("HandicapSet").Preload("Archers").Find(&tournaments)
	var handicapSets []models.HandicapSet
	db.Find(&handicapSets)
	var archers []models.Archer
	db.Preload("BowClass").Find(&archers)

	c.HTML(http.StatusOK, "tournaments", gin.H{
		"Title":        "Tournaments",
		"Tournaments":  tournaments,
		"HandicapSets": handicapSets,
		"Archers":      archers,
	})
}

func AddTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.PostForm("ID"))
	var tournament models.Tournament
	if id != 0 {
		if err := db.First(&tournament, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
			return
		}
	}

	if err := c.ShouldBind(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournament.TournamentType = c.PostForm("TournamentType")
	handicapSetID, _ := strconv.Atoi(c.PostForm("handicapSetID"))
	var handicapSet models.HandicapSet
	if err := db.First(&handicapSet, handicapSetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HandicapSet not found"})
		return
	}

	tournament.HandicapSetID = handicapSet.ID
	tournament.HandicapSet = handicapSet

	tournament.Date, _ = time.Parse("2006-01-02", c.PostForm("date"))
	archerIDs := c.PostFormArray("archers[]")

	if id == 0 {
		db.Create(&tournament)
	} else {
		db.Save(&tournament)
	}
	archerIDasInt := make([]uint, 0, len(archerIDs))
	for _, archerID := range archerIDs {
		id, _ := strconv.Atoi(archerID)
		archerIDasInt = append(archerIDasInt, uint(id))
	}
	tournamentHelper.UpsertTournamentArchers(db, tournament.ID, archerIDasInt)

	c.HTML(http.StatusOK, "tournament-oob", tournament)
	var handicapSets []models.HandicapSet
	db.Find(&handicapSets)
	var archers []models.Archer
	db.Preload("BowClass").Find(&archers)
	c.HTML(http.StatusOK, "tournament-form", gin.H{
		"HandicapSets": handicapSets,
		"Archers":      archers,
	})
}

func DeleteTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.Tournament{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	c.HTML(http.StatusNoContent, "", nil)
}
