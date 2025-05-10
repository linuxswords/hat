package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ShowBowClassesPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var bowClasses []models.BowClass
	db.Find(&bowClasses)
	c.HTML(http.StatusOK, "bowclasses.tmpl", gin.H{
		"Title":      "Bow Classes",
		"BowClasses": bowClasses,
	})
}

func AddBowClass(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var bowClass models.BowClass
	if err := c.ShouldBind(&bowClass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&bowClass)
	c.Redirect(http.StatusSeeOther, "/bowclasses")
}

func UpdateBowClass(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var bowClass models.BowClass
	if err := db.First(&bowClass, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BowClass not found"})
		return
	}
	if err := c.ShouldBindJSON(&bowClass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&bowClass)
	c.JSON(http.StatusOK, bowClass)
}

func DeleteBowClass(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.BowClass{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "BowClass not found"})
		return
	}
	c.Status(http.StatusOK)
}
