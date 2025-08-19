package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
)

// TournamentWithHandicap represents a tournament with its handicap set information
type TournamentWithHandicap struct {
	Tournament      models.Tournament
	HandicapSetName string
}

// TournamentsList handles GET /tournaments
func (h *TournamentHandlers) TournamentsList(c *gin.Context) {
	tournaments := h.TournamentRepo.GetAll()
	
	// Get current and upcoming tournaments for stats
	current := h.TournamentRepo.GetCurrent()
	upcoming := h.TournamentRepo.GetUpcoming()

	// Create tournaments with handicap set names
	var tournamentsWithHandicap []TournamentWithHandicap
	for _, tournament := range tournaments {
		tournamentData := TournamentWithHandicap{
			Tournament:      tournament,
			HandicapSetName: "",
		}
		
		// Get handicap set name if available
		if tournament.HandicapSetID != 0 {
			handicapSet, err := h.HandicapRepo.GetSetByID(tournament.HandicapSetID)
			if err == nil {
				tournamentData.HandicapSetName = handicapSet.Name
			}
		}
		
		tournamentsWithHandicap = append(tournamentsWithHandicap, tournamentData)
	}

	RenderWithLayout(c, "tournaments/index", gin.H{
		"title":          "Tournaments",
		"tournaments":    tournamentsWithHandicap,
		"currentCount":   len(current),
		"upcomingCount":  len(upcoming),
	})
}

// TournamentsNew handles GET /tournaments/new
func (h *TournamentHandlers) TournamentsNew(c *gin.Context) {
	// Get handicap sets for selection
	handicapSets := h.HandicapRepo.GetAllSets()

	RenderWithLayout(c, "tournaments/new", gin.H{
		"title":        "Create New Tournament",
		"handicapSets": handicapSets,
	})
}

// TournamentsCreate handles POST /tournaments
func (h *TournamentHandlers) TournamentsCreate(c *gin.Context) {
	name := c.PostForm("name")
	location := c.PostForm("location")
	startDateStr := c.PostForm("start_date")
	endDateStr := c.PostForm("end_date")
	handicapSetIDStr := c.PostForm("handicap_set_id")

	// Basic validation
	if name == "" || location == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, location, start date, and end date are required"})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Parse handicap set ID (optional)
	var handicapSetID uint
	if handicapSetIDStr != "" && handicapSetIDStr != "0" {
		id, err := strconv.Atoi(handicapSetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handicap set ID"})
			return
		}
		handicapSetID = uint(id)
	}

	// Create new tournament
	tournament := models.Tournament{
		Name:          name,
		Location:      location,
		StartDate:     startDate,
		EndDate:       endDate,
		HandicapSetID: handicapSetID,
	}

	// Add to repository
	_, err = h.TournamentRepo.Create(tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tournament"})
		return
	}

	// Redirect to tournament list
	c.Redirect(http.StatusFound, "/tournaments")
}

// TournamentsShow handles GET /tournaments/:id
func (h *TournamentHandlers) TournamentsShow(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Find tournament by ID
	tournament, err := h.TournamentRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get handicap set info if available
	var handicapSetName string
	if tournament.HandicapSetID != 0 {
		handicapSet, err := h.HandicapRepo.GetSetByID(tournament.HandicapSetID)
		if err == nil {
			handicapSetName = handicapSet.Name
		}
	}

	// Calculate tournament status
	now := time.Now()
	var status string
	var statusClass string
	
	if now.Before(tournament.StartDate) {
		status = "Upcoming"
		statusClass = "bg-blue-500"
	} else if now.After(tournament.EndDate) {
		status = "Completed"
		statusClass = "bg-gray-500"
	} else {
		status = "Active"
		statusClass = "bg-green-500"
	}

	RenderWithLayout(c, "tournaments/show", gin.H{
		"title":           "Tournament Details",
		"tournament":      tournament,
		"handicapSetName": handicapSetName,
		"status":          status,
		"statusClass":     statusClass,
	})
}

// TournamentsEdit handles GET /tournaments/:id/edit
func (h *TournamentHandlers) TournamentsEdit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Find tournament by ID
	tournament, err := h.TournamentRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get handicap sets for selection
	handicapSets := h.HandicapRepo.GetAllSets()

	RenderWithLayout(c, "tournaments/edit", gin.H{
		"title":        "Edit Tournament",
		"tournament":   tournament,
		"handicapSets": handicapSets,
	})
}

// TournamentsUpdate handles POST /tournaments/:id
func (h *TournamentHandlers) TournamentsUpdate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	name := c.PostForm("name")
	location := c.PostForm("location")
	startDateStr := c.PostForm("start_date")
	endDateStr := c.PostForm("end_date")
	handicapSetIDStr := c.PostForm("handicap_set_id")

	// Basic validation
	if name == "" || location == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, location, start date, and end date are required"})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Parse handicap set ID (optional)
	var handicapSetID uint
	if handicapSetIDStr != "" && handicapSetIDStr != "0" {
		hsID, err := strconv.Atoi(handicapSetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handicap set ID"})
			return
		}
		handicapSetID = uint(hsID)
	}

	// Update tournament
	updatedTournament := models.Tournament{
		Name:          name,
		Location:      location,
		StartDate:     startDate,
		EndDate:       endDate,
		HandicapSetID: handicapSetID,
	}

	_, err = h.TournamentRepo.Update(uint(id), updatedTournament)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Redirect to tournament details
	c.Redirect(http.StatusFound, "/tournaments/"+idParam)
}

// TournamentsDelete handles POST /tournaments/:id/delete
func (h *TournamentHandlers) TournamentsDelete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Delete tournament from repository
	err = h.TournamentRepo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Redirect to tournament list
	c.Redirect(http.StatusFound, "/tournaments")
}