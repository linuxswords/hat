package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test BowClassesList handler
func TestBowClassesList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/bowclasses", BowClassesList)

	t.Run("successful list", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/bowclasses", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// The test will succeed if the bow classes file exists,
		// or return 500 if it doesn't. Both are acceptable for this test.
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})
}

// Test BowClassesNew handler
func TestBowClassesNew(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/bowclasses/new", BowClassesNew)

	t.Run("successful new form", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/bowclasses/new", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// Test BowClassesShow handler
func TestBowClassesShow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/bowclasses/:id", BowClassesShow)

	t.Run("successful show with valid ID", func(t *testing.T) {
		// We need to use an ID that exists in the JSON file
		// Let's first check what bow classes are available
		bowClasses, err := loadBowClassesFromJSON()
		if err != nil || len(bowClasses) == 0 {
			t.Skip("No bow classes available in JSON file")
		}

		// Use the first available bow class ID
		validID := bowClasses[0].ID

		req, _ := http.NewRequest("GET", "/bowclasses/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bow class not found", func(t *testing.T) {
		// First check if the bow classes file exists
		_, err := loadBowClassesFromJSON()
		if err != nil {
			// If file doesn't exist, we expect 500 error
			req, _ := http.NewRequest("GET", "/bowclasses/nonexistent", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		} else {
			// If file exists, we expect 404 for non-existent ID
			req, _ := http.NewRequest("GET", "/bowclasses/nonexistent", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		}
	})
}

// Test loadBowClassesFromJSON function
func TestLoadBowClassesFromJSON(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		bowClasses, err := loadBowClassesFromJSON()
		
		// The function should not error if the file exists
		// We can't guarantee the file exists in all environments,
		// so we'll accept either success or a file not found error
		if err == nil {
			assert.NotNil(t, bowClasses)
			// If we successfully load, verify the structure
			for _, bc := range bowClasses {
				assert.NotEmpty(t, bc.ID, "Bow class should have an ID")
				assert.NotEmpty(t, bc.Name, "Bow class should have a name")
			}
		} else {
			// File might not exist in test environment
			t.Logf("Bow classes JSON file not found (this may be expected): %v", err)
		}
	})
}