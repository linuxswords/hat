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
	db.Preload("HandycapEntries.BowClass").Find(&handycapSets)
	c.HTML(http.StatusOK, "handycaps.tmpl", gin.H{
		"Title":        "Handycaps",
		"HandycapSets": handycapSets,
	})
}

func AddHandycapSet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var bowClasses []models.BowClass
	db.Find(&bowClasses)

	if c.Request.Method == http.MethodPost {
		var handycapSet models.HandycapSet
		if err := c.ShouldBind(&handycapSet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Bind handycap values
		for _, bowClass := range bowClasses {
			value, err := strconv.ParseFloat(c.PostForm("factor_"+strconv.Itoa(int(bowClass.ID))), 64)
			if err == nil {
				handycapSet.HandycapEntries = append(handycapSet.HandycapEntries, models.HandycapEntry{
					BowClassID: bowClass.ID,
					Value:      value,
				})
			}
		}

		db.Create(&handycapSet)
		c.Redirect(http.StatusSeeOther, "/handycaps")
		return
	}

	c.HTML(http.StatusOK, "add_handycapset.tmpl", gin.H{
		"Title":      "Add Handycap Set",
		"BowClasses": bowClasses,
	})
}

func EditHandycapSet(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var handycapSet models.HandycapSet
	if err := db.Preload("HandycapEntries.BowClass").First(&handycapSet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HandycapSet not found"})
		return
	}

	var bowClasses []models.BowClass
	db.Find(&bowClasses)

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&handycapSet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update handycap values
		for _, bowClass := range bowClasses {
			value, err := strconv.ParseFloat(c.PostForm("factor_"+strconv.Itoa(int(bowClass.ID))), 64)
			if err == nil {
				for i, entry := range handycapSet.HandycapEntries {
					if entry.BowClassID == bowClass.ID {
						handycapSet.HandycapEntries[i].Value = value
					}
				}
			}
		}

		db.Save(&handycapSet)
		c.Redirect(http.StatusSeeOther, "/handycaps")
		return
	}

	c.HTML(http.StatusOK, "edit_handycapset.tmpl", gin.H{
		"Title":       "Edit Handycap Set",
		"HandycapSet": handycapSet,
		"BowClasses":  bowClasses,
	})

	// var handycap models.HandycapEntry
	// if err := db.First(&handycap, id).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Handycap not found"})
	// 	return
	// }
	// if err := c.ShouldBindJSON(&handycap); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// db.Save(&handycap)
	// c.Status(http.StatusOK)
}

func DeleteHandycap(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.HandycapSet{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HandycapSet not found"})
		return
	}
	c.Status(http.StatusOK)
}
