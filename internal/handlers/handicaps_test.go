package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Setup test helpers for handicaps
func setupHandicapTest() (*gin.Engine, *HandicapHandlers, *MockHandicapRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create fresh mock
	mockRepo := new(MockHandicapRepository)
	
	// Create handlers with injected dependencies
	handlers := &HandicapHandlers{
		HandicapRepo: mockRepo,
	}
	
	return router, handlers, mockRepo
}

// Test HandicapsList handler
func TestHandicapsList(t *testing.T) {
	router, handlers, mockRepo := setupHandicapTest()

	router.GET("/handicaps", handlers.HandicapsList)

	t.Run("successful list", func(t *testing.T) {
		// Mock data
		mockHandicapSets := []models.HandicapSet{
			{Model: gorm.Model{ID: 1}, Name: "NFAA Indoor"},
			{Model: gorm.Model{ID: 2}, Name: "NFAA Outdoor"},
		}
		
		mockHandicaps1 := []models.Handicap{
			{Model: gorm.Model{ID: 1}, SetID: 1, BowClassID: "recurve", Factor: 1.0},
			{Model: gorm.Model{ID: 2}, SetID: 1, BowClassID: "compound", Factor: 0.8},
		}
		
		mockHandicaps2 := []models.Handicap{
			{Model: gorm.Model{ID: 3}, SetID: 2, BowClassID: "recurve", Factor: 1.1},
			{Model: gorm.Model{ID: 4}, SetID: 2, BowClassID: "compound", Factor: 0.9},
		}

		mockRepo.On("GetAllSets").Return(mockHandicapSets)
		mockRepo.On("GetHandicapsBySetID", uint(1)).Return(mockHandicaps1)
		mockRepo.On("GetHandicapsBySetID", uint(2)).Return(mockHandicaps2)

		req, _ := http.NewRequest("GET", "/handicaps", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

// Test HandicapsShow handler
func TestHandicapsShow(t *testing.T) {
	router, handlers, mockRepo := setupHandicapTest()

	router.GET("/handicaps/:id", handlers.HandicapsShow)

	t.Run("successful show", func(t *testing.T) {
		mockHandicapSet := &models.HandicapSet{
			Model: gorm.Model{ID: 1},
			Name: "NFAA Indoor",
		}
		
		mockHandicaps := []models.Handicap{
			{Model: gorm.Model{ID: 1}, SetID: 1, BowClassID: "recurve", Factor: 1.0},
			{Model: gorm.Model{ID: 2}, SetID: 1, BowClassID: "compound", Factor: 0.8},
			{Model: gorm.Model{ID: 3}, SetID: 1, BowClassID: "longbow", Factor: 1.2},
		}

		mockRepo.On("GetSetByID", uint(1)).Return(mockHandicapSet, nil)
		mockRepo.On("GetHandicapsBySetID", uint(1)).Return(mockHandicaps)

		req, _ := http.NewRequest("GET", "/handicaps/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/handicaps/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("handicap set not found", func(t *testing.T) {
		mockRepo.On("GetSetByID", uint(999)).Return(nil, errors.New("handicap set not found"))

		req, _ := http.NewRequest("GET", "/handicaps/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}