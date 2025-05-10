package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func ShowHandycapsPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var handycapSets []models.HandycapSet
	db.Preload("Handycaps").Find(&handycapSets)
	c.HTML(http.StatusOK, "handycaps.tmpl", gin.H{
		"Title":        "Handycaps",
		"HandycapSets": handycapSets,
	})
}

func AddHandycap(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var handycap models.Handycap
	if err := c.ShouldBind(&handycap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&handycap)
	c.Redirect(http.StatusSeeOther, "/handycaps")
}

func UpdateHandycap(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var handycap models.Handycap
	if err := db.First(&handycap, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Handycap not found"})
		return
	}
	if err := c.ShouldBindJSON(&handycap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&handycap)
	c.Status(http.StatusOK)
}

func DeleteHandycap(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.Handycap{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Handycap not found"})
		return
	}
	c.Status(http.StatusOK)
}
