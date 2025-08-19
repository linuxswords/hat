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

// Setup test helpers for tournaments
func setupTournamentTest() (*gin.Engine, *TournamentHandlers, *MockTournamentRepository, *MockHandicapRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create fresh mocks
	mockTournamentRepo := new(MockTournamentRepository)
	mockHandicapRepo := new(MockHandicapRepository)
	
	// Create handlers with injected dependencies
	handlers := &TournamentHandlers{
		TournamentRepo: mockTournamentRepo,
		HandicapRepo:   mockHandicapRepo,
	}
	
	return router, handlers, mockTournamentRepo, mockHandicapRepo
}

// Test TournamentsList handler
func TestTournamentsList(t *testing.T) {
	router, handlers, mockTournamentRepo, mockHandicapRepo := setupTournamentTest()

	router.GET("/tournaments", handlers.TournamentsList)

	t.Run("successful list", func(t *testing.T) {
		// Mock data
		now := time.Now()
		mockTournaments := []models.Tournament{
			{
				Model: gorm.Model{ID:            1},
				Name:          "Spring Championship",
				Location:      "Central Park",
				StartDate:     now.AddDate(0, 0, 10),
				EndDate:       now.AddDate(0, 0, 12),
				HandicapSetID: 1,
			},
			{
				Model: gorm.Model{ID:        2},
				Name:      "Summer Open",
				Location:  "Sport Center",
				StartDate: now.AddDate(0, 0, -5),
				EndDate:   now.AddDate(0, 0, -3),
			},
		}

		mockCurrent := []models.Tournament{}
		mockUpcoming := []models.Tournament{mockTournaments[0]}

		mockHandicapSet := &models.HandicapSet{
			Model: gorm.Model{ID:   1},
			Name: "NFAA Indoor",
		}

		mockTournamentRepo.On("GetAll").Return(mockTournaments)
		mockTournamentRepo.On("GetCurrent").Return(mockCurrent)
		mockTournamentRepo.On("GetUpcoming").Return(mockUpcoming)
		mockHandicapRepo.On("GetSetByID", uint(1)).Return(mockHandicapSet, nil)

		req, _ := http.NewRequest("GET", "/tournaments", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
	})
}

// Test TournamentsNew handler
func TestTournamentsNew(t *testing.T) {
	router, handlers, _, mockHandicapRepo := setupTournamentTest()

	router.GET("/tournaments/new", handlers.TournamentsNew)

	t.Run("successful new form", func(t *testing.T) {
		mockHandicapSets := []models.HandicapSet{
			{Model: gorm.Model{ID: 1}, Name: "NFAA Indoor"},
			{Model: gorm.Model{ID: 2}, Name: "NFAA Outdoor"},
		}

		mockHandicapRepo.On("GetAllSets").Return(mockHandicapSets)

		req, _ := http.NewRequest("GET", "/tournaments/new", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockHandicapRepo.AssertExpectations(t)
	})
}

// Test TournamentsCreate handler
func TestTournamentsCreate(t *testing.T) {
	router, handlers, mockTournamentRepo, _ := setupTournamentTest()

	router.POST("/tournaments", handlers.TournamentsCreate)

	t.Run("successful creation", func(t *testing.T) {
		startDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)

		expectedTournament := models.Tournament{
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     startDate,
			EndDate:       endDate,
			HandicapSetID: 1,
		}
		createdTournament := expectedTournament
		createdTournament.Model.ID = 1

		mockTournamentRepo.On("Create", expectedTournament).Return(&createdTournament, nil)

		formData := url.Values{}
		formData.Set("name", "Test Tournament")
		formData.Set("location", "Test Location")
		formData.Set("start_date", "2024-06-15")
		formData.Set("end_date", "2024-06-17")
		formData.Set("handicap_set_id", "1")

		req, _ := http.NewRequest("POST", "/tournaments", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/tournaments", w.Header().Get("Location"))
		mockTournamentRepo.AssertExpectations(t)
	})

	t.Run("missing required fields", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "Test Tournament")
		// Missing location, dates

		req, _ := http.NewRequest("POST", "/tournaments", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid start date", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "Test Tournament")
		formData.Set("location", "Test Location")
		formData.Set("start_date", "invalid-date")
		formData.Set("end_date", "2024-06-17")

		req, _ := http.NewRequest("POST", "/tournaments", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid end date", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "Test Tournament")
		formData.Set("location", "Test Location")
		formData.Set("start_date", "2024-06-15")
		formData.Set("end_date", "invalid-date")

		req, _ := http.NewRequest("POST", "/tournaments", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("end date before start date", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "Test Tournament")
		formData.Set("location", "Test Location")
		formData.Set("start_date", "2024-06-17")
		formData.Set("end_date", "2024-06-15")

		req, _ := http.NewRequest("POST", "/tournaments", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("repository error", func(t *testing.T) {
		startDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)

		expectedTournament := models.Tournament{
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     startDate,
			EndDate:       endDate,
			HandicapSetID: 0,
		}

		mockTournamentRepo.On("Create", expectedTournament).Return(nil, errors.New("database error"))

		formData := url.Values{}
		formData.Set("name", "Test Tournament")
		formData.Set("location", "Test Location")
		formData.Set("start_date", "2024-06-15")
		formData.Set("end_date", "2024-06-17")

		req, _ := http.NewRequest("POST", "/tournaments", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentsShow handler
func TestTournamentsShow(t *testing.T) {
	router, handlers, mockTournamentRepo, mockHandicapRepo := setupTournamentTest()

	router.GET("/tournaments/:id", handlers.TournamentsShow)

	t.Run("successful show with handicap set", func(t *testing.T) {
		now := time.Now()
		mockTournament := &models.Tournament{
			Model: gorm.Model{ID:            1},
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     now.AddDate(0, 0, 1),
			EndDate:       now.AddDate(0, 0, 3),
			HandicapSetID: 1,
		}

		mockHandicapSet := &models.HandicapSet{
			Model: gorm.Model{ID:   1},
			Name: "NFAA Indoor",
		}

		mockTournamentRepo.On("GetByID", uint(1)).Return(mockTournament, nil)
		mockHandicapRepo.On("GetSetByID", uint(1)).Return(mockHandicapSet, nil)

		req, _ := http.NewRequest("GET", "/tournaments/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
	})

	t.Run("successful show without handicap set", func(t *testing.T) {
		now := time.Now()
		mockTournament := &models.Tournament{
			Model: gorm.Model{ID:            1},
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     now.AddDate(0, 0, 1),
			EndDate:       now.AddDate(0, 0, 3),
			HandicapSetID: 0,
		}

		mockTournamentRepo.On("GetByID", uint(1)).Return(mockTournament, nil)

		req, _ := http.NewRequest("GET", "/tournaments/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", uint(999)).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentsEdit handler
func TestTournamentsEdit(t *testing.T) {
	router, handlers, mockTournamentRepo, mockHandicapRepo := setupTournamentTest()

	router.GET("/tournaments/:id/edit", handlers.TournamentsEdit)

	t.Run("successful edit form", func(t *testing.T) {
		now := time.Now()
		mockTournament := &models.Tournament{
			Model: gorm.Model{ID:            1},
			Name:          "Test Tournament",
			Location:      "Test Location",
			StartDate:     now.AddDate(0, 0, 1),
			EndDate:       now.AddDate(0, 0, 3),
			HandicapSetID: 1,
		}

		mockHandicapSets := []models.HandicapSet{
			{Model: gorm.Model{ID: 1}, Name: "NFAA Indoor"},
			{Model: gorm.Model{ID: 2}, Name: "NFAA Outdoor"},
		}

		mockTournamentRepo.On("GetByID", uint(1)).Return(mockTournament, nil)
		mockHandicapRepo.On("GetAllSets").Return(mockHandicapSets)

		req, _ := http.NewRequest("GET", "/tournaments/1/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTournamentRepo.AssertExpectations(t)
		mockHandicapRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tournaments/invalid/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("GetByID", uint(999)).Return(nil, errors.New("tournament not found"))

		req, _ := http.NewRequest("GET", "/tournaments/999/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentsUpdate handler
func TestTournamentsUpdate(t *testing.T) {
	router, handlers, mockTournamentRepo, _ := setupTournamentTest()

	router.POST("/tournaments/:id", handlers.TournamentsUpdate)

	t.Run("successful update", func(t *testing.T) {
		startDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)

		updatedTournament := models.Tournament{
			Name:          "Updated Tournament",
			Location:      "Updated Location",
			StartDate:     startDate,
			EndDate:       endDate,
			HandicapSetID: 2,
		}
		returnTournament := updatedTournament
		returnTournament.Model.ID = 1

		mockTournamentRepo.On("Update", uint(1), updatedTournament).Return(&returnTournament, nil)

		formData := url.Values{}
		formData.Set("name", "Updated Tournament")
		formData.Set("location", "Updated Location")
		formData.Set("start_date", "2024-06-15")
		formData.Set("end_date", "2024-06-17")
		formData.Set("handicap_set_id", "2")

		req, _ := http.NewRequest("POST", "/tournaments/1", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/tournaments/1", w.Header().Get("Location"))
		mockTournamentRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "Updated Tournament")
		formData.Set("location", "Updated Location")
		formData.Set("start_date", "2024-06-15")
		formData.Set("end_date", "2024-06-17")

		req, _ := http.NewRequest("POST", "/tournaments/invalid", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		startDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)

		updatedTournament := models.Tournament{
			Name:          "Updated Tournament",
			Location:      "Updated Location",
			StartDate:     startDate,
			EndDate:       endDate,
			HandicapSetID: 0,
		}

		mockTournamentRepo.On("Update", uint(999), updatedTournament).Return(nil, errors.New("tournament not found"))

		formData := url.Values{}
		formData.Set("name", "Updated Tournament")
		formData.Set("location", "Updated Location")
		formData.Set("start_date", "2024-06-15")
		formData.Set("end_date", "2024-06-17")

		req, _ := http.NewRequest("POST", "/tournaments/999", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}

// Test TournamentsDelete handler
func TestTournamentsDelete(t *testing.T) {
	router, handlers, mockTournamentRepo, _ := setupTournamentTest()

	router.POST("/tournaments/:id/delete", handlers.TournamentsDelete)

	t.Run("successful delete", func(t *testing.T) {
		mockTournamentRepo.On("Delete", uint(1)).Return(nil)

		req, _ := http.NewRequest("POST", "/tournaments/1/delete", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/tournaments", w.Header().Get("Location"))
		mockTournamentRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/tournaments/invalid/delete", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tournament not found", func(t *testing.T) {
		mockTournamentRepo.On("Delete", uint(999)).Return(errors.New("tournament not found"))

		req, _ := http.NewRequest("POST", "/tournaments/999/delete", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockTournamentRepo.AssertExpectations(t)
	})
}