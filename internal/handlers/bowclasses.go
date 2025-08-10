package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type BowClass struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// loadBowClassesFromJSON loads bow classes from the JSON file
func loadBowClassesFromJSON() ([]BowClass, error) {
	data, err := os.ReadFile("doc/data/bowclasses/bowclasses.json")
	if err != nil {
		return nil, err
	}

	var bowClasses []BowClass
	err = json.Unmarshal(data, &bowClasses)
	if err != nil {
		return nil, err
	}

	return bowClasses, nil
}

// BowClassesList handles GET /bowclasses
func BowClassesList(c *gin.Context) {
	bowClasses, err := loadBowClassesFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load bow classes"})
		return
	}

	RenderWithLayout(c, "bowclasses/index", gin.H{
		"title":      "Bow Classes",
		"bowClasses": bowClasses,
	})
}

// BowClassesNew handles GET /bowclasses/new
func BowClassesNew(c *gin.Context) {
	RenderWithLayout(c, "bowclasses/new", gin.H{
		"title": "Create New Bow Class",
	})
}

// BowClassesShow handles GET /bowclasses/:id
func BowClassesShow(c *gin.Context) {
	id := c.Param("id")

	bowClasses, err := loadBowClassesFromJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load bow classes"})
		return
	}

	var bowClass *BowClass
	for _, bc := range bowClasses {
		if bc.ID == id {
			bowClass = &bc
			break
		}
	}

	if bowClass == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bow class not found"})
		return
	}

	RenderWithLayout(c, "bowclasses/show", gin.H{
		"title":    "Bow Class Details",
		"bowClass": bowClass,
	})
}
