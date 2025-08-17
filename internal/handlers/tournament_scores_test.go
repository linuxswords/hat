package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Setup test helpers for tournament scores
func setupTournamentScoresTest() (*gin.Engine, *TournamentScoreHandlers, *MockScoreRepository, *MockTournamentRepository, *MockArcherRepository, *MockHandicapRepository, *MockTournamentParticipationRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create fresh mocks
	mockScoreRepo := new(MockScoreRepository)
	mockTournamentRepo := new(MockTournamentRepository)
	mockArcherRepo := new(MockArcherRepository)
	mockHandicapRepo := new(MockHandicapRepository)
	mockParticipationRepo := new(MockTournamentParticipationRepository)
	
	// Create handlers with injected dependencies
	handlers := &TournamentScoreHandlers{
		ScoreRepo:         mockScoreRepo,
		TournamentRepo:    mockTournamentRepo,
		ArcherRepo:        mockArcherRepo,
		HandicapRepo:      mockHandicapRepo,
		ParticipationRepo: mockParticipationRepo,
	}
	
	return router, handlers, mockScoreRepo, mockTournamentRepo, mockArcherRepo, mockHandicapRepo, mockParticipationRepo
}

// Test TournamentScoresIndex handler
func TestTournamentScoresIndex(t *testing.T) {
	router, handlers, mockScoreRepo, mockTournamentRepo, mockArcherRepo, mockHandicapRepo, mockParticipationRepo := setupTournamentScoresTest()

	router.GET("/tournaments/:id/scores", handlers.TournamentScoresIndex)

	t.Run("successful scores index with handicap", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			ID:            1,
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 1),
			HandicapSetID: 1,
		}

		// Mock participants
		participations := []models.TournamentParticipation{
			{TournamentID: 1, ArcherID: 1},
			{TournamentID: 1, ArcherID: 2},
		}

		// Mock archers
		archer1 := &models.Archer{ID: 1, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}
		archer2 := &models.Archer{ID: 2, Name: "Jane Smith", BowClass: "compound", Gender: "Female", Email: "jane@example.com"}

		// Mock scores
		adjustedScore1 := 285.5
		scores := []models.Score{
			{ID: 1, ArcherID: 1, TournamentID: 1, RawScore: 280, AdjustedScore: &adjustedScore1},
		}

		// Mock handicap set and handicaps
		handicapSet := &models.HandicapSet{ID: 1, Name: "NFAA Indoor"}
		handicap1 := &models.Handicap{Factor: 1.02}
		handicap2 := &models.Handicap{Factor: 0.95}

		// Set up expectations
		mockTournamentRepo.On("GetByID", 1).Return(tournament, nil)
		mockParticipationRepo.On("GetByTournamentID", 1).Return(participations)
		mockScoreRepo.On("GetByTournamentID", 1).Return(scores)
		mockHandicapRepo.On("GetSetByID", 1).Return(handicapSet, nil)
		mockArcherRepo.On("GetByID", 1).Return(archer1, nil)
		mockArcherRepo.On("GetByID", 2).Return(archer2, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "recurve").Return(handicap1, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "compound").Return(handicap2, nil)

		req, _ := http.NewRequest("GET", "/tournaments/1/scores", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockParticipationRepo.AssertExpectations(t)
		mockScoreRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid/scores", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", 999).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999/scores", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentScoresEdit handler
func TestTournamentScoresEdit(t *testing.T) {
	router, handlers, mockScoreRepo, mockTournamentRepo, mockArcherRepo, mockHandicapRepo, mockParticipationRepo := setupTournamentScoresTest()

	router.GET("/tournaments/:id/scores/edit", handlers.TournamentScoresEdit)

	t.Run("successful scores edit", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			ID:            1,
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 1),
			HandicapSetID: 1,
		}

		// Mock participants
		participations := []models.TournamentParticipation{
			{TournamentID: 1, ArcherID: 1},
		}

		// Mock archer
		archer1 := &models.Archer{ID: 1, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}

		// Mock existing score
		adjustedScore := 285.5
		scores := []models.Score{
			{ID: 1, ArcherID: 1, TournamentID: 1, RawScore: 280, AdjustedScore: &adjustedScore},
		}

		// Mock handicap set and handicap
		handicapSet := &models.HandicapSet{ID: 1, Name: "NFAA Indoor"}
		handicap := &models.Handicap{Factor: 1.02}

		// Set up expectations
		mockTournamentRepo.On("GetByID", 1).Return(tournament, nil)
		mockParticipationRepo.On("GetByTournamentID", 1).Return(participations)
		mockScoreRepo.On("GetByTournamentID", 1).Return(scores)
		mockHandicapRepo.On("GetSetByID", 1).Return(handicapSet, nil)
		mockArcherRepo.On("GetByID", 1).Return(archer1, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "recurve").Return(handicap, nil)

		req, _ := http.NewRequest("GET", "/tournaments/1/scores/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockParticipationRepo.AssertExpectations(t)
		mockScoreRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid/scores/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", 999).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999/scores/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentScoresUpdate handler
func TestTournamentScoresUpdate(t *testing.T) {
	router, handlers, mockScoreRepo, mockTournamentRepo, mockArcherRepo, mockHandicapRepo, _ := setupTournamentScoresTest()

	router.POST("/tournaments/:id/scores", handlers.TournamentScoresUpdate)

	t.Run("successful scores update - create new score", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			ID:            1,
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 1),
			HandicapSetID: 1,
		}

		// Mock archer
		archer := &models.Archer{ID: 1, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}

		// Mock handicap
		handicap := &models.Handicap{Factor: 1.02}

		// Expected score
		adjustedScore := 285.6
		expectedScore := models.Score{
			ArcherID:      1,
			TournamentID:  1,
			BowClassID:    "recurve",
			RawScore:      280,
			AdjustedScore: &adjustedScore,
			EnteredBy:     "tournament_organizer",
		}
		createdScore := expectedScore
		createdScore.ID = 1

		// Set up expectations
		mockTournamentRepo.On("GetByID", 1).Return(tournament, nil)
		mockArcherRepo.On("GetByID", 1).Return(archer, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "recurve").Return(handicap, nil)
		mockScoreRepo.On("CalculateAdjustedScore", 280, 1.02).Return(285.6)
		mockScoreRepo.On("GetByTournamentAndArcher", 1, 1).Return(nil, errors.New("not found"))
		mockScoreRepo.On("Create", mock.MatchedBy(func(score models.Score) bool {
			return score.ArcherID == 1 && score.TournamentID == 1 && score.RawScore == 280
		})).Return(&createdScore, nil)

		formData := url.Values{}
		formData.Set("score_1", "280")

		req, _ := http.NewRequest("POST", "/tournaments/1/scores", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/scores?success=1")
		mockTournamentRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
		mockScoreRepo.AssertExpectations(t)
	})

}

// Test TournamentScoresUpdate - update existing score (separate test function)
func TestTournamentScoresUpdate_UpdateExistingScore(t *testing.T) {
	router, handlers, mockScoreRepo, mockTournamentRepo, mockArcherRepo, mockHandicapRepo, _ := setupTournamentScoresTest()

	router.POST("/tournaments/:id/scores", handlers.TournamentScoresUpdate)

	t.Run("successful scores update - update existing score", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			ID:            1,
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 1),
			HandicapSetID: 1,
		}

		// Mock archer
		archer := &models.Archer{ID: 1, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}

		// Mock handicap
		handicap := &models.Handicap{Factor: 1.02}

		// Existing score
		existingScore := &models.Score{ID: 1, ArcherID: 1, TournamentID: 1, RawScore: 275}

		// Expected updated score
		adjustedScore := 285.6
		expectedScore := models.Score{
			ArcherID:      1,
			TournamentID:  1,
			BowClassID:    "recurve",
			RawScore:      280,
			AdjustedScore: &adjustedScore,
			EnteredBy:     "tournament_organizer",
		}
		updatedScore := expectedScore
		updatedScore.ID = 1

		// Set up expectations
		mockTournamentRepo.On("GetByID", 1).Return(tournament, nil)
		mockArcherRepo.On("GetByID", 1).Return(archer, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "recurve").Return(handicap, nil)
		mockScoreRepo.On("CalculateAdjustedScore", 280, 1.02).Return(285.6)
		mockScoreRepo.On("GetByTournamentAndArcher", 1, 1).Return(existingScore, nil)
		mockScoreRepo.On("UpdateByTournamentAndArcher", 1, 1, mock.MatchedBy(func(score models.Score) bool {
			return score.ArcherID == 1 && score.TournamentID == 1 && score.RawScore == 280
		})).Return(&updatedScore, nil)

		formData := url.Values{}
		formData.Set("score_1", "280")

		req, _ := http.NewRequest("POST", "/tournaments/1/scores", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/scores?success=1")
		mockTournamentRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
		mockScoreRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("score_1", "280")

		req, _ := http.NewRequest("POST", "/tournaments/invalid/scores", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", 999).Return(nil, errors.New("tournament not found"))

		formData := url.Values{}
		formData.Set("score_1", "280")

		req, _ := http.NewRequest("POST", "/tournaments/999/scores", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})

	t.Run("invalid score value", func(t *testing.T) {
		tournament := &models.Tournament{
			ID:            1,
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 1),
			HandicapSetID: 0,
		}

		mockTournamentRepo.On("GetByID", 1).Return(tournament, nil)

		formData := url.Values{}
		formData.Set("score_1", "invalid")

		req, _ := http.NewRequest("POST", "/tournaments/1/scores", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/scores/edit?error=save_failed")
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentScoresRankings handler
func TestTournamentScoresRankings(t *testing.T) {
	router, handlers, mockScoreRepo, mockTournamentRepo, mockArcherRepo, mockHandicapRepo, _ := setupTournamentScoresTest()

	router.GET("/tournaments/:id/rankings", handlers.TournamentScoresRankings)

	t.Run("successful rankings", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			ID:            1,
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 1),
			HandicapSetID: 1,
		}

		// Mock sorted scores
		adjustedScore1 := 290.0
		adjustedScore2 := 285.5
		scores := []models.Score{
			{ID: 1, ArcherID: 1, TournamentID: 1, RawScore: 285, AdjustedScore: &adjustedScore1},
			{ID: 2, ArcherID: 2, TournamentID: 1, RawScore: 280, AdjustedScore: &adjustedScore2},
		}

		// Mock archers
		archer1 := &models.Archer{ID: 1, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}
		archer2 := &models.Archer{ID: 2, Name: "Jane Smith", BowClass: "compound", Gender: "Female", Email: "jane@example.com"}

		// Mock handicap set and handicaps
		handicapSet := &models.HandicapSet{ID: 1, Name: "NFAA Indoor"}
		handicap1 := &models.Handicap{Factor: 1.02}
		handicap2 := &models.Handicap{Factor: 0.95}

		// Set up expectations
		mockTournamentRepo.On("GetByID", 1).Return(tournament, nil)
		mockScoreRepo.On("GetByTournamentIDSorted", 1).Return(scores)
		mockHandicapRepo.On("GetSetByID", 1).Return(handicapSet, nil)
		mockArcherRepo.On("GetByID", 1).Return(archer1, nil)
		mockArcherRepo.On("GetByID", 2).Return(archer2, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "recurve").Return(handicap1, nil)
		mockHandicapRepo.On("GetHandicapByBowClass", 1, "compound").Return(handicap2, nil)

		req, _ := http.NewRequest("GET", "/tournaments/1/rankings", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockScoreRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid/rankings", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", 999).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999/rankings", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}