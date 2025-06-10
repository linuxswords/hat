package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ShowArchersPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var archers []models.Archer
	db.Preload("BowClass").Find(&archers)
	var bowClasses []models.BowClass
	db.Find(&bowClasses)
	c.HTML(http.StatusOK, "archers", gin.H{
		"Title":      "Archers",
		"Archers":    archers,
		"BowClasses": bowClasses,
	})
}

func AddArcher(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.PostForm("ID"))
	var archer models.Archer
	if id != 0 {
		if err := db.First(&archer, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
			return
		}
	}

	if err := c.ShouldBind(&archer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bowClass models.BowClass
	if err := db.First(&bowClass, archer.BowClassID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	archer.BowClass = bowClass

	if id == 0 {
		db.Create(&archer)
	} else {
		db.Save(&archer)
	}

	c.HTML(http.StatusOK, "archerForm", archer)
	c.HTML(http.StatusCreated, "archer-oob", archer)
}

func UpdateArcher(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var archer models.Archer
	if err := db.First(&archer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}
	if err := c.ShouldBindJSON(&archer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&archer)
	c.JSON(http.StatusOK, archer)
}

func GetArcher(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var archer models.Archer
	if err := db.Preload("BowClass").First(&archer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}
	c.JSON(http.StatusOK, archer)
}

func DeleteArcher(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.Archer{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/archers")
}
