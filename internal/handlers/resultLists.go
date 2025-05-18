package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func ShowResultLists(c *gin.Context) {
	// Fetch tournaments from the database
	var scores []models.Score
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&scores)

	var archers []models.Archer
	db.Preload("Score").Find(&archers)

	c.HTML(http.StatusOK, "resultLists.tmpl", gin.H{
		"Title":   "HAT - Result Lists",
		"Content": "Result Lists",
		"Scores":  scores,
	})
}
