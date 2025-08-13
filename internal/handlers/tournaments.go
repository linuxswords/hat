package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/linuxswords/hat/internal/repositories"
)

// TournamentRepository defines the contract for tournament data operations
type TournamentRepository interface {
	GetAll() []models.Tournament
	GetByID(id int) (*models.Tournament, error)
	Create(tournament models.Tournament) (*models.Tournament, error)
	Update(id int, tournament models.Tournament) (*models.Tournament, error)
	Delete(id int) error
	GetUpcoming() []models.Tournament
	GetCurrent() []models.Tournament
	GetByHandicapSet(handicapSetID int) []models.Tournament
}

var tournamentRepo TournamentRepository = repositories.NewTournamentRepository()

// TournamentWithHandicap represents a tournament with its handicap set information
type TournamentWithHandicap struct {
	Tournament      models.Tournament
	HandicapSetName string
}

// TournamentsList handles GET /tournaments
func TournamentsList(c *gin.Context) {
	tournaments := tournamentRepo.GetAll()
	
	// Get current and upcoming tournaments for stats
	current := tournamentRepo.GetCurrent()
	upcoming := tournamentRepo.GetUpcoming()

	// Create tournaments with handicap set names
	var tournamentsWithHandicap []TournamentWithHandicap
	for _, tournament := range tournaments {
		tournamentData := TournamentWithHandicap{
			Tournament:      tournament,
			HandicapSetName: "",
		}
		
		// Get handicap set name if available
		if tournament.HandicapSetID != nil {
			handicapSet, err := handicapRepo.GetSetByID(*tournament.HandicapSetID)
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
func TournamentsNew(c *gin.Context) {
	// Get handicap sets for selection
	handicapSets := handicapRepo.GetAllSets()

	RenderWithLayout(c, "tournaments/new", gin.H{
		"title":        "Create New Tournament",
		"handicapSets": handicapSets,
	})
}

// TournamentsCreate handles POST /tournaments
func TournamentsCreate(c *gin.Context) {
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
	var handicapSetID *int
	if handicapSetIDStr != "" && handicapSetIDStr != "0" {
		id, err := strconv.Atoi(handicapSetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handicap set ID"})
			return
		}
		handicapSetID = &id
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
	_, err = tournamentRepo.Create(tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tournament"})
		return
	}

	// Redirect to tournament list
	c.Redirect(http.StatusFound, "/tournaments")
}

// TournamentsShow handles GET /tournaments/:id
func TournamentsShow(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Find tournament by ID
	tournament, err := tournamentRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get handicap set info if available
	var handicapSetName string
	if tournament.HandicapSetID != nil {
		handicapSet, err := handicapRepo.GetSetByID(*tournament.HandicapSetID)
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
func TournamentsEdit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Find tournament by ID
	tournament, err := tournamentRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get handicap sets for selection
	handicapSets := handicapRepo.GetAllSets()

	RenderWithLayout(c, "tournaments/edit", gin.H{
		"title":        "Edit Tournament",
		"tournament":   tournament,
		"handicapSets": handicapSets,
	})
}

// TournamentsUpdate handles POST /tournaments/:id
func TournamentsUpdate(c *gin.Context) {
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
	var handicapSetID *int
	if handicapSetIDStr != "" && handicapSetIDStr != "0" {
		hsID, err := strconv.Atoi(handicapSetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handicap set ID"})
			return
		}
		handicapSetID = &hsID
	}

	// Update tournament
	updatedTournament := models.Tournament{
		Name:          name,
		Location:      location,
		StartDate:     startDate,
		EndDate:       endDate,
		HandicapSetID: handicapSetID,
	}

	_, err = tournamentRepo.Update(id, updatedTournament)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Redirect to tournament details
	c.Redirect(http.StatusFound, "/tournaments/"+idParam)
}

// TournamentsDelete handles POST /tournaments/:id/delete
func TournamentsDelete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Delete tournament from repository
	err = tournamentRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Redirect to tournament list
	c.Redirect(http.StatusFound, "/tournaments")
}