package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
)

// HandicapSetData represents a handicap set with its handicaps for template rendering
type HandicapSetData struct {
	Set       models.HandicapSet
	Handicaps []models.Handicap
}

// HandicapsList handles GET /handicaps
func (h *HandicapHandlers) HandicapsList(c *gin.Context) {
	handicapSets := h.HandicapRepo.GetAllSets()
	
	// Get handicaps for each set
	var setsData []HandicapSetData
	for _, set := range handicapSets {
		handicaps := h.HandicapRepo.GetHandicapsBySetID(set.ID)
		// Sort handicaps by bow class ID for consistent display
		sort.Slice(handicaps, func(i, j int) bool {
			return handicaps[i].BowClassID < handicaps[j].BowClassID
		})
		
		setsData = append(setsData, HandicapSetData{
			Set:       set,
			Handicaps: handicaps,
		})
	}

	RenderWithLayout(c, "handicaps/index", gin.H{
		"title":        "Handicap Sets",
		"handicapSets": setsData,
	})
}

// HandicapsShow handles GET /handicaps/:id
func (h *HandicapHandlers) HandicapsShow(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handicap set ID"})
		return
	}

	// Find handicap set by ID
	handicapSet, err := h.HandicapRepo.GetSetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Handicap set not found"})
		return
	}

	// Get all handicaps for this set
	handicaps := h.HandicapRepo.GetHandicapsBySetID(uint(id))
	
	// Sort handicaps by bow class ID for consistent display
	sort.Slice(handicaps, func(i, j int) bool {
		return handicaps[i].BowClassID < handicaps[j].BowClassID
	})

	// Load bow classes for context
	bowClasses, err := loadBowClassesFromJSON()
	if err != nil {
		// Continue without bow class names if we can't load them
		bowClasses = []BowClass{}
	}

	// Create a map for quick bow class name lookup
	bowClassNames := make(map[string]string)
	for _, bc := range bowClasses {
		bowClassNames[bc.ID] = bc.Name
	}

	RenderWithLayout(c, "handicaps/show", gin.H{
		"title":          "Handicap Set Details",
		"handicapSet":    handicapSet,
		"handicaps":      handicaps,
		"bowClassNames":  bowClassNames,
	})
}