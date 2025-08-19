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
	"gorm.io/gorm"
)

// Setup test helpers for tournament archers
func setupTournamentArcherTest() (*gin.Engine, *TournamentArcherHandlers, *MockTournamentRepository, *MockArcherRepository, *MockTournamentParticipationRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create fresh mocks
	mockTournamentRepo := new(MockTournamentRepository)
	mockArcherRepo := new(MockArcherRepository)
	mockParticipationRepo := new(MockTournamentParticipationRepository)
	
	// Create handlers with injected dependencies
	handlers := &TournamentArcherHandlers{
		TournamentRepo:    mockTournamentRepo,
		ArcherRepo:        mockArcherRepo,
		ParticipationRepo: mockParticipationRepo,
	}
	
	return router, handlers, mockTournamentRepo, mockArcherRepo, mockParticipationRepo
}

// Test TournamentArchersIndex handler
func TestTournamentArchersIndex(t *testing.T) {
	router, handlers, mockTournamentRepo, mockArcherRepo, mockParticipationRepo := setupTournamentArcherTest()

	router.GET("/tournaments/:id/archers", handlers.TournamentArchersIndex)

	t.Run("successful index", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			Model:     gorm.Model{ID: 1},
			Name:      "Test Tournament",
			Location:  "Test Location",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 1),
		}

		// Mock participations
		participations := []models.TournamentParticipation{
			{Model: gorm.Model{ID: 1}, TournamentID: 1, ArcherID: 1, Status: "registered"},
			{Model: gorm.Model{ID: 2}, TournamentID: 1, ArcherID: 2, Status: "checked_in"},
		}

		// Mock archers
		archer1 := &models.Archer{Model: gorm.Model{ID: 1}, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}
		archer2 := &models.Archer{Model: gorm.Model{ID: 2}, Name: "Jane Smith", BowClass: "compound", Gender: "Female", Email: "jane@example.com"}

		// Set up expectations
		mockTournamentRepo.On("GetByID", uint(1)).Return(tournament, nil)
		mockParticipationRepo.On("GetByTournamentID", uint(1)).Return(participations)
		mockArcherRepo.On("GetByID", uint(1)).Return(archer1, nil)
		mockArcherRepo.On("GetByID", uint(2)).Return(archer2, nil)

		req, _ := http.NewRequest("GET", "/tournaments/1/archers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockParticipationRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid/archers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", uint(999)).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999/archers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})

	t.Run("handles missing archer gracefully", func(t *testing.T) {
		// Create separate setup for this test to avoid mock conflicts
		subRouter, subHandlers, subMockTournamentRepo, subMockArcherRepo, subMockParticipationRepo := setupTournamentArcherTest()
		subRouter.GET("/tournaments/:id/archers", subHandlers.TournamentArchersIndex)

		// Mock tournament
		tournament := &models.Tournament{
			Model:     gorm.Model{ID: 1},
			Name:      "Test Tournament",
			Location:  "Test Location",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 1),
		}

		// Mock participations with one valid and one invalid archer
		participations := []models.TournamentParticipation{
			{Model: gorm.Model{ID: 1}, TournamentID: 1, ArcherID: 1, Status: "registered"},
			{Model: gorm.Model{ID: 2}, TournamentID: 1, ArcherID: 999, Status: "registered"}, // This archer doesn't exist
		}

		archer1 := &models.Archer{Model: gorm.Model{ID: 1}, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"}

		subMockTournamentRepo.On("GetByID", uint(1)).Return(tournament, nil)
		subMockParticipationRepo.On("GetByTournamentID", uint(1)).Return(participations)
		subMockArcherRepo.On("GetByID", uint(1)).Return(archer1, nil)
		subMockArcherRepo.On("GetByID", uint(999)).Return(nil, errors.New("archer not found"))

		req, _ := http.NewRequest("GET", "/tournaments/1/archers", nil)
		w := httptest.NewRecorder()
		subRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		subMockTournamentRepo.AssertExpectations(t)
		subMockParticipationRepo.AssertExpectations(t)
		subMockArcherRepo.AssertExpectations(t)
	})
}

// Test TournamentArchersAdd handler
func TestTournamentArchersAdd(t *testing.T) {
	router, handlers, mockTournamentRepo, mockArcherRepo, mockParticipationRepo := setupTournamentArcherTest()

	router.GET("/tournaments/:id/archers/add", handlers.TournamentArchersAdd)

	t.Run("successful add form", func(t *testing.T) {
		// Mock tournament
		tournament := &models.Tournament{
			Model:     gorm.Model{ID: 1},
			Name:      "Test Tournament",
			Location:  "Test Location",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 1),
		}

		// Mock all archers
		allArchers := []models.Archer{
			{Model: gorm.Model{ID: 1}, Name: "John Doe", BowClass: "recurve", Gender: "Male", Email: "john@example.com"},
			{Model: gorm.Model{ID: 2}, Name: "Jane Smith", BowClass: "compound", Gender: "Female", Email: "jane@example.com"},
			{Model: gorm.Model{ID: 3}, Name: "Bob Wilson", BowClass: "longbow", Gender: "Male", Email: "bob@example.com"},
		}

		// Mock existing participations (archer 1 is already registered)
		participations := []models.TournamentParticipation{
			{Model: gorm.Model{ID: 1}, TournamentID: 1, ArcherID: 1, Status: "registered"},
		}

		mockTournamentRepo.On("GetByID", uint(1)).Return(tournament, nil)
		mockArcherRepo.On("GetAll").Return(allArchers)
		mockParticipationRepo.On("GetByTournamentID", uint(1)).Return(participations)

		req, _ := http.NewRequest("GET", "/tournaments/1/archers/add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockArcherRepo.AssertExpectations(t)
		mockParticipationRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid/archers/add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", uint(999)).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999/archers/add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentArchersCreate handler
func TestTournamentArchersCreate(t *testing.T) {
	router, handlers, _, _, mockParticipationRepo := setupTournamentArcherTest()

	router.POST("/tournaments/:id/archers", handlers.TournamentArchersCreate)

	t.Run("successful creation with single archer", func(t *testing.T) {
		expectedParticipation := models.TournamentParticipation{
			TournamentID: 1,
			ArcherID:     2,
			Status:       "registered",
		}
		createdParticipation := expectedParticipation
		createdParticipation.ID = 1

		mockParticipationRepo.On("Create", expectedParticipation).Return(&createdParticipation, nil)

		formData := url.Values{}
		formData.Set("archer_ids", "2")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/archers?success=1")
		mockParticipationRepo.AssertExpectations(t)
	})

	t.Run("successful creation with multiple archers", func(t *testing.T) {
		expectedParticipation1 := models.TournamentParticipation{
			TournamentID: 1,
			ArcherID:     2,
			Status:       "registered",
		}
		expectedParticipation2 := models.TournamentParticipation{
			TournamentID: 1,
			ArcherID:     3,
			Status:       "registered",
		}
		
		createdParticipation1 := expectedParticipation1
		createdParticipation1.ID = 1
		createdParticipation2 := expectedParticipation2
		createdParticipation2.ID = 2

		mockParticipationRepo.On("Create", expectedParticipation1).Return(&createdParticipation1, nil)
		mockParticipationRepo.On("Create", expectedParticipation2).Return(&createdParticipation2, nil)

		formData := url.Values{}
		formData.Add("archer_ids", "2")
		formData.Add("archer_ids", "3")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/archers?success=2")
		mockParticipationRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("archer_ids", "2")

		req, _ := http.NewRequest("POST", "/tournaments/invalid/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("no archers selected", func(t *testing.T) {
		formData := url.Values{}
		// No archer_ids provided

		req, _ := http.NewRequest("POST", "/tournaments/1/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid archer ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("archer_ids", "invalid")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/archers/add?error=registration_failed")
	})

	t.Run("repository error - all fail", func(t *testing.T) {
		// Create separate setup for this test to avoid mock conflicts
		subRouter, subHandlers, _, _, subMockParticipationRepo := setupTournamentArcherTest()
		subRouter.POST("/tournaments/:id/archers", subHandlers.TournamentArchersCreate)

		expectedParticipation := models.TournamentParticipation{
			TournamentID: 1,
			ArcherID:     2,
			Status:       "registered",
		}

		subMockParticipationRepo.On("Create", expectedParticipation).Return(nil, errors.New("database error"))

		formData := url.Values{}
		formData.Set("archer_ids", "2")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		subRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/tournaments/1/archers/add?error=registration_failed")
		subMockParticipationRepo.AssertExpectations(t)
	})
}

// Test TournamentArchersRemove handler
func TestTournamentArchersRemove(t *testing.T) {
	router, handlers, _, _, mockParticipationRepo := setupTournamentArcherTest()

	router.POST("/tournaments/:id/archers/:archer_id/remove", handlers.TournamentArchersRemove)

	t.Run("successful removal", func(t *testing.T) {
		mockParticipationRepo.On("DeleteByTournamentAndArcher", uint(1), uint(2)).Return(nil)

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/2/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/tournaments/1/archers", w.Header().Get("Location"))
		mockParticipationRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/tournaments/invalid/archers/2/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid archer ID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/tournaments/1/archers/invalid/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("participation not found", func(t *testing.T) {
		mockParticipationRepo.On("DeleteByTournamentAndArcher", uint(1), uint(999)).Return(errors.New("participation not found"))

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/999/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockParticipationRepo.AssertExpectations(t)
	})
}

// Test TournamentArchersUpdateStatus handler
func TestTournamentArchersUpdateStatus(t *testing.T) {
	router, handlers, _, _, mockParticipationRepo := setupTournamentArcherTest()

	router.POST("/tournaments/:id/archers/:archer_id/status", handlers.TournamentArchersUpdateStatus)

	t.Run("successful status update", func(t *testing.T) {
		existingParticipation := &models.TournamentParticipation{
			Model:        gorm.Model{ID: 1},
			TournamentID: 1,
			ArcherID:     2,
			Status:       "registered",
		}

		updatedParticipation := models.TournamentParticipation{
			Model:        gorm.Model{ID: 1},
			TournamentID: 1,
			ArcherID:     2,
			Status:       "checked_in",
		}

		mockParticipationRepo.On("GetByTournamentAndArcher", uint(1), uint(2)).Return(existingParticipation, nil)
		mockParticipationRepo.On("Update", uint(1), updatedParticipation).Return(&updatedParticipation, nil)

		formData := url.Values{}
		formData.Set("status", "checked_in")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/2/status", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/tournaments/1/archers", w.Header().Get("Location"))
		mockParticipationRepo.AssertExpectations(t)
	})

	t.Run("invalid tournament ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("status", "checked_in")

		req, _ := http.NewRequest("POST", "/tournaments/invalid/archers/2/status", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid archer ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("status", "checked_in")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/invalid/status", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid status", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("status", "invalid_status")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/2/status", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("participation not found", func(t *testing.T) {
		mockParticipationRepo.On("GetByTournamentAndArcher", uint(1), uint(999)).Return(nil, errors.New("participation not found"))

		formData := url.Values{}
		formData.Set("status", "checked_in")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/999/status", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockParticipationRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		// Create separate setup for this test to avoid mock conflicts
		subRouter, subHandlers, _, _, subMockParticipationRepo := setupTournamentArcherTest()
		subRouter.POST("/tournaments/:id/archers/:archer_id/status", subHandlers.TournamentArchersUpdateStatus)

		existingParticipation := &models.TournamentParticipation{
			Model:        gorm.Model{ID: 1},
			TournamentID: 1,
			ArcherID:     2,
			Status:       "registered",
		}

		updatedParticipation := models.TournamentParticipation{
			Model:        gorm.Model{ID: 1},
			TournamentID: 1,
			ArcherID:     2,
			Status:       "checked_in",
		}

		subMockParticipationRepo.On("GetByTournamentAndArcher", uint(1), uint(2)).Return(existingParticipation, nil)
		subMockParticipationRepo.On("Update", uint(1), updatedParticipation).Return(nil, errors.New("database error"))

		formData := url.Values{}
		formData.Set("status", "checked_in")

		req, _ := http.NewRequest("POST", "/tournaments/1/archers/2/status", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		subRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		subMockParticipationRepo.AssertExpectations(t)
	})
}