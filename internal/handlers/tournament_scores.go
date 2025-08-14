package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/linuxswords/hat/internal/repositories"
)

// ScoreRepository defines the contract for score operations
type ScoreRepository interface {
	GetAll() []models.Score
	GetByTournamentID(tournamentID int) []models.Score
	GetByTournamentIDSorted(tournamentID int) []models.Score
	GetByArcherID(archerID int) []models.Score
	GetByTournamentAndArcher(tournamentID int, archerID int) (*models.Score, error)
	Create(score models.Score) (*models.Score, error)
	Update(id int, score models.Score) (*models.Score, error)
	UpdateByTournamentAndArcher(tournamentID int, archerID int, score models.Score) (*models.Score, error)
	Delete(id int) error
	DeleteByTournamentAndArcher(tournamentID int, archerID int) error
	GetScoreCountByTournament(tournamentID int) int
	CalculateAdjustedScore(rawScore int, handicapFactor float64) float64
}

var scoreRepo ScoreRepository = repositories.NewScoreRepository()

// Repository instances (these may be duplicated across files, but Go will handle this)
var tournamentScoresRepo TournamentRepository = repositories.NewTournamentRepository()
var handicapScoresRepo HandicapRepository = repositories.NewHandicapRepository()
var archerScoresRepo ArcherRepository = repositories.NewArcherRepository()
var participationScoresRepo TournamentParticipationRepository = repositories.NewTournamentParticipationRepository()

// TournamentScoresIndex handles GET /tournaments/:id/scores
func TournamentScoresIndex(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get tournament details
	tournament, err := tournamentScoresRepo.GetByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get tournament participants
	participations := participationScoresRepo.GetByTournamentID(tournamentID)
	participantMap := make(map[int]models.TournamentParticipation)
	for _, p := range participations {
		participantMap[p.ArcherID] = p
	}

	// Get scores for this tournament
	scores := scoreRepo.GetByTournamentID(tournamentID)
	scoreMap := make(map[int]models.Score)
	for _, s := range scores {
		scoreMap[s.ArcherID] = s
	}

	// Get handicap set information
	var handicapSet *models.HandicapSet
	if tournament.HandicapSetID != 0 {
		handicapSet, _ = handicapScoresRepo.GetSetByID(tournament.HandicapSetID)
	}

	// Create tournament score views for all participants
	var tournamentScores []models.TournamentScoreView
	for _, participation := range participations {
		archer, err := archerScoresRepo.GetByID(participation.ArcherID)
		if err != nil {
			continue
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

		// Get handicap factor if handicap set is used
		var handicapFactor float64
		if handicapSet != nil {
			handicap, err := handicapScoresRepo.GetHandicapByBowClass(tournament.HandicapSetID, archer.BowClass)
			if err == nil {
				handicapFactor = handicap.Factor
			}
		}

		// Create score view
		scoreView := models.TournamentScoreView{
			Archer:         *archer,
			BowClassName:   bowClassName,
			HandicapFactor: handicapFactor,
		}

		// Add score if exists
		if score, exists := scoreMap[archer.ID]; exists {
			scoreView.Score = score
		}

		tournamentScores = append(tournamentScores, scoreView)
	}

	// Sort by name for consistent display
	sort.Slice(tournamentScores, func(i, j int) bool {
		return tournamentScores[i].Archer.Name < tournamentScores[j].Archer.Name
	})

	// Calculate statistics
	totalScores := len(scores)
	var avgRawScore float64
	var avgAdjustedScore float64
	if totalScores > 0 {
		var rawSum, adjustedSum float64
		adjustedCount := 0
		for _, score := range scores {
			rawSum += float64(score.RawScore)
			if score.AdjustedScore != nil {
				adjustedSum += *score.AdjustedScore
				adjustedCount++
			}
		}
		avgRawScore = rawSum / float64(totalScores)
		if adjustedCount > 0 {
			avgAdjustedScore = adjustedSum / float64(adjustedCount)
		}
	}

	// Calculate progress percentage
	var progressPercent float64
	if len(participations) > 0 {
		progressPercent = float64(totalScores*100) / float64(len(participations))
	}

	RenderWithLayout(c, "tournaments/scores", gin.H{
		"title":             "Tournament Scores",
		"tournament":        tournament,
		"handicapSet":       handicapSet,
		"tournamentScores":  tournamentScores,
		"totalParticipants": len(participations),
		"totalScores":       totalScores,
		"avgRawScore":       avgRawScore,
		"avgAdjustedScore":  avgAdjustedScore,
		"progressPercent":   progressPercent,
	})
}

// TournamentScoresEdit handles GET /tournaments/:id/scores/edit
func TournamentScoresEdit(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get tournament details
	tournament, err := tournamentScoresRepo.GetByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get tournament participants
	participations := participationScoresRepo.GetByTournamentID(tournamentID)

	// Get scores for this tournament
	scores := scoreRepo.GetByTournamentID(tournamentID)
	scoreMap := make(map[int]models.Score)
	for _, s := range scores {
		scoreMap[s.ArcherID] = s
	}

	// Get handicap set information
	var handicapSet *models.HandicapSet
	if tournament.HandicapSetID != 0 {
		handicapSet, _ = handicapScoresRepo.GetSetByID(tournament.HandicapSetID)
	}

	// Create editable score entries
	var editableScores []models.TournamentScoreView
	for _, participation := range participations {
		archer, err := archerScoresRepo.GetByID(participation.ArcherID)
		if err != nil {
			continue
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

		// Get handicap factor if handicap set is used
		var handicapFactor float64
		if handicapSet != nil {
			handicap, err := handicapScoresRepo.GetHandicapByBowClass(tournament.HandicapSetID, archer.BowClass)
			if err == nil {
				handicapFactor = handicap.Factor
			}
		}

		scoreView := models.TournamentScoreView{
			Archer:         *archer,
			BowClassName:   bowClassName,
			HandicapFactor: handicapFactor,
		}

		// Add existing score if available
		if score, exists := scoreMap[archer.ID]; exists {
			scoreView.Score = score
		}

		editableScores = append(editableScores, scoreView)
	}

	// Sort by name for consistent display
	sort.Slice(editableScores, func(i, j int) bool {
		return editableScores[i].Archer.Name < editableScores[j].Archer.Name
	})

	RenderWithLayout(c, "tournaments/edit_scores", gin.H{
		"title":          "Edit Tournament Scores",
		"tournament":     tournament,
		"handicapSet":    handicapSet,
		"editableScores": editableScores,
	})
}

// TournamentScoresUpdate handles POST /tournaments/:id/scores
func TournamentScoresUpdate(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get tournament details for handicap calculations
	tournament, err := tournamentScoresRepo.GetByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Parse form data
	err = c.Request.ParseForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Process each archer's score from form data
	form := c.Request.Form
	var errors []string
	var successCount int


	for key, values := range form {
		// Skip non-score fields
		if len(values) == 0 || values[0] == "" {
			continue
		}

		// Parse archer ID from field name (format: "score_<archer_id>")
		if len(key) <= 6 || key[:6] != "score_" {
			continue
		}

		archerIDStr := key[6:]
		archerID, err := strconv.Atoi(archerIDStr)
		if err != nil {
			errors = append(errors, "Invalid archer ID: "+archerIDStr)
			continue
		}

		rawScore, err := strconv.Atoi(values[0])
		if err != nil || rawScore < 0 {
			errors = append(errors, "Invalid score for archer ID "+archerIDStr)
			continue
		}

		// Get archer for bow class
		archer, err := archerScoresRepo.GetByID(archerID)
		if err != nil {
			errors = append(errors, "Archer not found: "+archerIDStr)
			continue
		}

		// Create or update score
		score := models.Score{
			ArcherID:     archerID,
			TournamentID: tournamentID,
			BowClassID:   archer.BowClass,
			RawScore:     rawScore,
			EnteredBy:    "tournament_organizer", // Could be enhanced with user management
		}

		// Calculate adjusted score if handicap set is used
		if tournament.HandicapSetID != 0 {
			handicap, err := handicapScoresRepo.GetHandicapByBowClass(tournament.HandicapSetID, archer.BowClass)
			if err == nil {
				adjustedScore := scoreRepo.CalculateAdjustedScore(rawScore, handicap.Factor)
				score.AdjustedScore = &adjustedScore
			}
		}

		// Try to update existing score first
		_, err = scoreRepo.GetByTournamentAndArcher(tournamentID, archerID)
		if err == nil {
			// Update existing score
			_, err = scoreRepo.UpdateByTournamentAndArcher(tournamentID, archerID, score)
			if err != nil {
				errors = append(errors, "Failed to update score for archer ID "+archerIDStr+": "+err.Error())
			} else {
				successCount++
			}
		} else {
			// Create new score
			_, err = scoreRepo.Create(score)
			if err != nil {
				errors = append(errors, "Failed to create score for archer ID "+archerIDStr+": "+err.Error())
			} else {
				successCount++
			}
		}
	}

	// Redirect with status
	if successCount > 0 {
		c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/scores?success="+strconv.Itoa(successCount))
	} else if len(errors) > 0 {
		c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/scores/edit?error=save_failed")
	} else {
		c.Redirect(http.StatusFound, "/tournaments/"+tournamentIDParam+"/scores")
	}
}

// TournamentScoresRankings handles GET /tournaments/:id/rankings
func TournamentScoresRankings(c *gin.Context) {
	tournamentIDParam := c.Param("id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}

	// Get tournament details
	tournament, err := tournamentScoresRepo.GetByID(tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournament not found"})
		return
	}

	// Get sorted scores for ranking
	scores := scoreRepo.GetByTournamentIDSorted(tournamentID)

	// Get handicap set information
	var handicapSet *models.HandicapSet
	if tournament.HandicapSetID != 0 {
		handicapSet, _ = handicapScoresRepo.GetSetByID(tournament.HandicapSetID)
	}

	// Create ranking list
	var rankings []models.TournamentScoreView
	for rank, score := range scores {
		archer, err := archerScoresRepo.GetByID(score.ArcherID)
		if err != nil {
			continue
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

		// Get handicap factor
		var handicapFactor float64
		if handicapSet != nil {
			handicap, err := handicapScoresRepo.GetHandicapByBowClass(tournament.HandicapSetID, archer.BowClass)
			if err == nil {
				handicapFactor = handicap.Factor
			}
		}

		ranking := models.TournamentScoreView{
			Score:          score,
			Archer:         *archer,
			BowClassName:   bowClassName,
			HandicapFactor: handicapFactor,
			Rank:           rank + 1,
		}

		rankings = append(rankings, ranking)
	}

	RenderWithLayout(c, "tournaments/rankings", gin.H{
		"title":       "Tournament Rankings",
		"tournament":  tournament,
		"handicapSet": handicapSet,
		"rankings":    rankings,
	})
}

