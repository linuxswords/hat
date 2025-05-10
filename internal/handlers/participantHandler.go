package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ShowParticipantsPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var participants []models.Participant
	db.Preload("BowClass").Find(&participants)
	var bowClasses []models.BowClass
	db.Find(&bowClasses)
	c.HTML(http.StatusOK, "participants.tmpl", gin.H{
		"Title":        "Participants",
		"Participants": participants,
		"BowClasses":   bowClasses,
	})
}

func AddParticipant(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var participant models.Participant
	if err := c.ShouldBind(&participant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&participant)
	c.Redirect(http.StatusSeeOther, "/participants")
}

func UpdateParticipant(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var participant models.Participant
	if err := db.First(&participant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Participant not found"})
		return
	}
	if err := c.ShouldBindJSON(&participant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&participant)
	c.JSON(http.StatusOK, participant)
}

func DeleteParticipant(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.Participant{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Participant not found"})
		return
	}
	c.Status(http.StatusOK)
}
