package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func ShowHandicapsPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var handicapSets []models.HandicapSet
	db.Preload("HandicapEntries.BowClass").Find(&handicapSets)
	c.HTML(http.StatusOK, "handicaps", gin.H{
		"Title":        "Handicaps",
		"HandicapSets": handicapSets,
	})
}

func AddHandicapSet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var bowClasses []models.BowClass
	db.Find(&bowClasses)

	if c.Request.Method == http.MethodPost {
		var handicapSet models.HandicapSet
		if err := c.ShouldBind(&handicapSet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Bind handicap values
		for _, bowClass := range bowClasses {
			value, err := strconv.ParseFloat(c.PostForm("factor_"+strconv.Itoa(int(bowClass.ID))), 64)
			if err == nil {
				handicapSet.HandicapEntries = append(handicapSet.HandicapEntries, models.HandicapEntry{
					BowClassID: bowClass.ID,
					Value:      value,
				})
			}
		}

		db.Create(&handicapSet)
		c.Redirect(http.StatusSeeOther, "/handicaps")
		return
	}

	c.HTML(http.StatusOK, "addHandicapset", gin.H{
		"Title":      "Add Handicap Set",
		"BowClasses": bowClasses,
	})
}

func EditHandicapSet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var handicapSet models.HandicapSet
	if err := db.Preload("HandicapEntries.BowClass").First(&handicapSet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HandicapSet not found"})
		return
	}

	var bowClasses []models.BowClass
	db.Find(&bowClasses)

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&handicapSet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update handicap values
		for _, bowClass := range bowClasses {
			value, err := strconv.ParseFloat(c.PostForm("factor_"+strconv.Itoa(int(bowClass.ID))), 64)
			if err == nil {
				for i, entry := range handicapSet.HandicapEntries {
					if entry.BowClassID == bowClass.ID {
						handicapSet.HandicapEntries[i].Value = value
					}
				}
			}
		}

		db.Save(&handicapSet)
		c.Redirect(http.StatusSeeOther, "/handicaps")
		return
	}

	c.HTML(http.StatusOK, "editHandicapset", gin.H{
		"Title":       "Edit Handicap Set",
		"HandicapSet": handicapSet,
		"BowClasses":  bowClasses,
	})

	// var handicap models.HandicapEntry
	// if err := db.First(&handicap, id).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Handicap not found"})
	// 	return
	// }
	// if err := c.ShouldBindJSON(&handicap); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// db.Save(&handicap)
	// c.Status(http.StatusOK)
}

func DeleteHandicap(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.HandicapSet{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HandicapSet not found"})
		return
	}
	c.Status(http.StatusOK)
}
