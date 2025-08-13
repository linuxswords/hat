package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/linuxswords/hat/internal/repositories"
)

// HandicapRepository defines the contract for handicap data operations
type HandicapRepository interface {
	GetAllSets() []models.HandicapSet
	GetSetByID(id int) (*models.HandicapSet, error)
	GetHandicapsBySetID(setID int) []models.Handicap
	GetHandicapByBowClass(setID int, bowClassID string) (*models.Handicap, error)
	GetActiveSet() (*models.HandicapSet, error)
}

var handicapRepo HandicapRepository = repositories.NewHandicapRepository()

// HandicapSetData represents a handicap set with its handicaps for template rendering
type HandicapSetData struct {
	Set       models.HandicapSet
	Handicaps []models.Handicap
}

// HandicapsList handles GET /handicaps
func HandicapsList(c *gin.Context) {
	handicapSets := handicapRepo.GetAllSets()
	
	// Get handicaps for each set
	var setsData []HandicapSetData
	for _, set := range handicapSets {
		handicaps := handicapRepo.GetHandicapsBySetID(set.ID)
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
func HandicapsShow(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handicap set ID"})
		return
	}

	// Find handicap set by ID
	handicapSet, err := handicapRepo.GetSetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Handicap set not found"})
		return
	}

	// Get all handicaps for this set
	handicaps := handicapRepo.GetHandicapsBySetID(id)
	
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