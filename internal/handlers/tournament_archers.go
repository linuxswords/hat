package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
)

// TournamentArchersIndex handles GET /tournaments/:id/archers
func (h *TournamentArcherHandlers) TournamentArchersIndex(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get tournament details
	tournament, err := h.TournamentRepo.GetByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get participations for this tournament
	participations := h.ParticipationRepo.GetByTournamentID(tournamentID)

	// Get archer details and bow class names
	var tournamentArchers []models.TournamentArcherView
	for _, participation := range participations {
		archer, err := h.ArcherRepo.GetByID(participation.ArcherID)
		if err != nil {
			continue // Skip if archer not found
		}

		// Get bow class name
		bowClasses, err := loadBowClassesFromJSON()
		var bowClassName string
		if err == nil {
			for _, bc := range bowClasses {
				if bc.ID == archer.BowClass {
					bowClassName = bc.Name
					break
				}
			}
		}

		tournamentArcher := models.TournamentArcherView{
			Participation: participation,
			Archer:        *archer,
			BowClassName:  bowClassName,
		}
		tournamentArchers = append(tournamentArchers, tournamentArcher)
	}

	// Sort by archer name
	sort.Slice(tournamentArchers, func(i, j int) bool {
		return tournamentArchers[i].Archer.Name < tournamentArchers[j].Archer.Name
	})

	// Get bow class statistics
	bowClassStats := make(map[string]int)
	for _, ta := range tournamentArchers {
		bowClassStats[ta.Archer.BowClass]++
	}

	RenderWithLayout(c, "tournaments/archers", gin.H{
		"title":            "Tournament Archers",
		"tournament":       tournament,
		"tournamentArchers": tournamentArchers,
		"totalArchers":     len(tournamentArchers),
		"bowClassStats":    bowClassStats,
	})
}

// TournamentArchersAdd handles GET /tournaments/:id/archers/add
func (h *TournamentArcherHandlers) TournamentArchersAdd(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get tournament details
	tournament, err := h.TournamentRepo.GetByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get all archers
	allArchers := h.ArcherRepo.GetAll()

	// Get already registered archer IDs
	participations := h.ParticipationRepo.GetByTournamentID(tournamentID)
	registeredArcherIDs := make(map[int]bool)
	for _, participation := range participations {
		registeredArcherIDs[participation.ArcherID] = true
	}

	// Filter out already registered archers
	var availableArchers []models.Archer
	for _, archer := range allArchers {
		if !registeredArcherIDs[archer.ID] {
			availableArchers = append(availableArchers, archer)
		}
	}

	// Sort available archers by name
	sort.Slice(availableArchers, func(i, j int) bool {
		return availableArchers[i].Name < availableArchers[j].Name
	})

	// Get bow classes for display
	bowClasses, _ := loadBowClassesFromJSON()
	bowClassNames := make(map[string]string)
	for _, bc := range bowClasses {
		bowClassNames[bc.ID] = bc.Name
	}

	RenderWithLayout(c, "tournaments/add_archers", gin.H{
		"title":            "Add Archers to Tournament",
		"tournament":       tournament,
		"availableArchers": availableArchers,
		"bowClassNames":    bowClassNames,
	})
}

// TournamentArchersCreate handles POST /tournaments/:id/archers
func (h *TournamentArcherHandlers) TournamentArchersCreate(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get archer IDs from form (can be multiple)
	archerIDsStr := c.PostFormArray("archer_ids")
	if len(archerIDsStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No archers selected"})
		return
	}

	var successCount int
	var errors []string

	// Add each archer to the tournament
	for _, archerIDStr := range archerIDsStr {
		archerID, err := strconv.Atoi(archerIDStr)
		if err != nil {
			errors = append(errors, "Invalid archer ID: "+archerIDStr)
			continue
		}

		// Create participation
		participation := models.TournamentParticipation{
			TournamentID: tournamentID,
			ArcherID:     archerID,
			Status:       "registered",
		}

		_, err = h.ParticipationRepo.Create(participation)
		if err != nil {
			errors = append(errors, "Failed to register archer ID "+archerIDStr+": "+err.Error())
		} else {
			successCount++
		}
	}

	// Redirect back to tournament archers page with status
	if successCount > 0 {
		c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/archers?success="+strconv.Itoa(successCount))
	} else {
		c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/archers/add?error=registration_failed")
	}
}

// TournamentArchersRemove handles POST /tournaments/:id/archers/:archer_id/remove
func (h *TournamentArcherHandlers) TournamentArchersRemove(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	archerIDParam := c.Param("archer_id")

	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	archerID, err := strconv.Atoi(archerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid archer ID"})
		return
	}

	// Remove participation
	err = h.ParticipationRepo.DeleteByTournamentAndArcher(tournamentID, archerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Participation not found"})
		return
	}

	// Redirect back to tournament archers page
	c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/archers")
}

// TournamentArchersUpdateStatus handles POST /tournaments/:id/archers/:archer_id/status
func (h *TournamentArcherHandlers) TournamentArchersUpdateStatus(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	archerIDParam := c.Param("archer_id")

	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	archerID, err := strconv.Atoi(archerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid archer ID"})
		return
	}

	newStatus := c.PostForm("status")
	validStatuses := map[string]bool{
		"registered":  true,
		"checked_in":  true,
		"competing":   true,
		"completed":   true,
		"withdrawn":   true,
	}

	if !validStatuses[newStatus] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Get existing participation
	participation, err := h.ParticipationRepo.GetByTournamentAndArcher(tournamentID, archerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Participation not found"})
		return
	}

	// Update status
	participation.Status = newStatus
	_, err = h.ParticipationRepo.Update(participation.ID, *participation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	// Redirect back to tournament archers page
	c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/archers")
}