package handlers

import (
	"log"
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
	db.Preload("HandicapSet").Preload("Archers").Find(&tournaments)
	var handicapSets []models.HandicapSet
	db.Find(&handicapSets)
	var archers []models.Archer
	db.Preload("BowClass").Find(&archers)

	c.HTML(http.StatusOK, "tournaments.tmpl", gin.H{
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
	for _, archerID := range archerIDs {
		id, _ := strconv.Atoi(archerID)
		var archer models.Archer
		if err := db.First(&archer, id).Error; err == nil {
			db.Model(&tournament).Association("Archers").Append(&archer)
		}
	}

	if id == 0 {
		db.Create(&tournament)
	} else {
		db.Save(&tournament)
	}

	c.Redirect(http.StatusSeeOther, "/tournaments")
}

func GetTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var tournament models.Tournament
	if err := db.Preload("Archers").First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	if err := db.Preload("BowClass").Find(&tournament.Archers).Error; err != nil {
		log.Panic("Could not load bowclasses for the archers", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "bowclasses for archers not found"})
		return
	}
	c.JSON(http.StatusOK, tournament)
}

func UpdateTournament(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var tournament models.Tournament
	if err := db.First(&tournament, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}
	if err := c.ShouldBind(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tournament.Date, _ = time.Parse("2006-01-02", c.PostForm("date"))
	tournament.TournamentType = c.PostForm("TournamentType")
	handicapSetID, _ := strconv.Atoi(c.PostForm("handicapSetID"))
	tournament.HandicapSetID = uint(handicapSetID)
	// Update archers associated with the tournament
	archerIDs := c.PostFormArray("archers[]")
	var archers []models.Archer
	for _, archerID := range archerIDs {
		id, _ := strconv.Atoi(archerID)
		var archer models.Archer
		if err := db.First(&archer, id).Error; err == nil {
			archers = append(archers, archer)
		}
	}
	db.Model(&tournament).Association("Archers").Replace(archers)

	db.Save(&tournament)
	c.Redirect(http.StatusSeeOther, "/tournaments")
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
