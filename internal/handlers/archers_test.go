package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// setupArcherTest creates a fresh test environment for each test
func setupArcherTest() (*gin.Engine, *ArcherHandlers, *MockArcherRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create fresh mock
	mockRepo := new(MockArcherRepository)
	
	// Create handlers with injected dependencies
	handlers := &ArcherHandlers{
		ArcherRepo: mockRepo,
	}
	
	return router, handlers, mockRepo
}

// Test ArchersList handler
func TestArchersList(t *testing.T) {
	router, handlers, mockRepo := setupArcherTest()

	// Setup route
	router.GET("/archers", handlers.ArchersList)

	t.Run("successful list", func(t *testing.T) {
		// Mock data
		mockArchers := []models.Archer{
			{Model: gorm.Model{ID: 1}, Name: "John Doe", Gender: "Male", BowClass: "recurve", Email: "john@example.com"},
			{Model: gorm.Model{ID: 2}, Name: "Jane Smith", Gender: "Female", BowClass: "compound", Email: "jane@example.com"},
		}

		mockRepo.On("GetAll").Return(mockArchers)

		// Make request
		req, _ := http.NewRequest("GET", "/archers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

// Test ArchersNew handler
func TestArchersNew(t *testing.T) {
	router, handlers, _ := setupArcherTest()

	router.GET("/archers/new", handlers.ArchersNew)

	t.Run("successful new form", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/archers/new", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Note: Testing the error case where loadBowClassesForArchers() fails
	// would require mocking the file system or testing in an environment
	// where the bow classes file doesn't exist. Since the function loads
	// from a static file, we'd need to refactor it to be injectable
	// to properly test the error case.
}

// Test ArchersCreate handler
func TestArchersCreate(t *testing.T) {
	router, handlers, mockRepo := setupArcherTest()

	router.POST("/archers", handlers.ArchersCreate)

	t.Run("successful creation", func(t *testing.T) {
		// Setup mock
		expectedArcher := models.Archer{
			Name:     "John Doe",
			Gender:   "Male",
			BowClass: "recurve",
			Email:    "john@example.com",
		}
		createdArcher := expectedArcher
		createdArcher.ID = 1

		mockRepo.On("Create", expectedArcher).Return(&createdArcher, nil)

		// Prepare form data
		formData := url.Values{}
		formData.Set("name", "John Doe")
		formData.Set("gender", "Male")
		formData.Set("bow_class", "recurve")
		formData.Set("email", "john@example.com")

		req, _ := http.NewRequest("POST", "/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/archers", w.Header().Get("Location"))
		mockRepo.AssertExpectations(t)
	})

	t.Run("missing required fields", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "John Doe")
		// Missing other required fields

		req, _ := http.NewRequest("POST", "/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("repository error", func(t *testing.T) {
		// Create fresh setup for this subtest
		subRouter, subHandlers, subMockRepo := setupArcherTest()
		subRouter.POST("/archers", subHandlers.ArchersCreate)
		
		expectedArcher := models.Archer{
			Name:     "John Doe",
			Gender:   "Male",
			BowClass: "recurve",
			Email:    "john@example.com",
		}

		subMockRepo.On("Create", expectedArcher).Return(nil, errors.New("database error"))

		formData := url.Values{}
		formData.Set("name", "John Doe")
		formData.Set("gender", "Male")
		formData.Set("bow_class", "recurve")
		formData.Set("email", "john@example.com")

		req, _ := http.NewRequest("POST", "/archers", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		subRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		subMockRepo.AssertExpectations(t)
	})
}

// Test ArchersShow handler
func TestArchersShow(t *testing.T) {
	router, handlers, mockRepo := setupArcherTest()

	router.GET("/archers/:id", handlers.ArchersShow)

	t.Run("successful show", func(t *testing.T) {
		mockArcher := &models.Archer{
			Model:    gorm.Model{ID: 1},
			Name:     "John Doe",
			Gender:   "Male",
			BowClass: "recurve",
			Email:    "john@example.com",
		}

		mockRepo.On("GetByID", uint(1)).Return(mockArcher, nil)

		req, _ := http.NewRequest("GET", "/archers/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/archers/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("archer not found", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("archer not found"))

		req, _ := http.NewRequest("GET", "/archers/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

// Test ArchersEdit handler
func TestArchersEdit(t *testing.T) {
	router, handlers, mockRepo := setupArcherTest()

	router.GET("/archers/:id/edit", handlers.ArchersEdit)

	t.Run("successful edit form", func(t *testing.T) {
		mockArcher := &models.Archer{
			Model:    gorm.Model{ID: 1},
			Name:     "John Doe",
			Gender:   "Male",
			BowClass: "recurve",
			Email:    "john@example.com",
		}

		mockRepo.On("GetByID", uint(1)).Return(mockArcher, nil)

		req, _ := http.NewRequest("GET", "/archers/1/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/archers/invalid/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("archer not found", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("archer not found"))

		req, _ := http.NewRequest("GET", "/archers/999/edit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

// Test ArchersUpdate handler
func TestArchersUpdate(t *testing.T) {
	router, handlers, mockRepo := setupArcherTest()

	router.POST("/archers/:id", handlers.ArchersUpdate)

	t.Run("successful update", func(t *testing.T) {
		updatedArcher := models.Archer{
			Name:     "John Doe Updated",
			Gender:   "Male",
			BowClass: "compound",
			Email:    "john.updated@example.com",
		}
		returnArcher := updatedArcher
		returnArcher.ID = 1

		mockRepo.On("Update", uint(1), updatedArcher).Return(&returnArcher, nil)

		formData := url.Values{}
		formData.Set("name", "John Doe Updated")
		formData.Set("gender", "Male")
		formData.Set("bow_class", "compound")
		formData.Set("email", "john.updated@example.com")

		req, _ := http.NewRequest("POST", "/archers/1", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/archers/1", w.Header().Get("Location"))
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "John Doe")
		formData.Set("gender", "Male")
		formData.Set("bow_class", "recurve")
		formData.Set("email", "john@example.com")

		req, _ := http.NewRequest("POST", "/archers/invalid", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "John Doe")
		// Missing other required fields

		req, _ := http.NewRequest("POST", "/archers/1", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("archer not found", func(t *testing.T) {
		updatedArcher := models.Archer{
			Name:     "John Doe",
			Gender:   "Male",
			BowClass: "recurve",
			Email:    "john@example.com",
		}

		mockRepo.On("Update", uint(999), updatedArcher).Return(nil, errors.New("archer not found"))

		formData := url.Values{}
		formData.Set("name", "John Doe")
		formData.Set("gender", "Male")
		formData.Set("bow_class", "recurve")
		formData.Set("email", "john@example.com")

		req, _ := http.NewRequest("POST", "/archers/999", strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}

// Test ArchersDelete handler
func TestArchersDelete(t *testing.T) {
	router, handlers, mockRepo := setupArcherTest()

	router.POST("/archers/:id/delete", handlers.ArchersDelete)

	t.Run("successful delete", func(t *testing.T) {
		mockRepo.On("Delete", uint(1)).Return(nil)

		req, _ := http.NewRequest("POST", "/archers/1/delete", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "/archers", w.Header().Get("Location"))
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/archers/invalid/delete", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("archer not found", func(t *testing.T) {
		mockRepo.On("Delete", uint(999)).Return(errors.New("archer not found"))

		req, _ := http.NewRequest("POST", "/archers/999/delete", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}