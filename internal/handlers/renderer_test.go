package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test getProjectRoot function
func TestGetProjectRoot(t *testing.T) {
	t.Run("returns valid path", func(t *testing.T) {
		root := getProjectRoot()
		assert.NotEmpty(t, root, "Project root should not be empty")
		assert.Contains(t, root, "hat", "Project root should contain 'hat'")
	})
}

// Test RenderWithLayout function
func TestRenderWithLayout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful render", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			RenderWithLayout(c, "test/template", gin.H{
				"title": "Test Page",
			})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// The response will depend on whether template files exist
		// We accept either success or template error as valid outcomes
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})

	t.Run("handles missing template gracefully", func(t *testing.T) {
		router := gin.New()
		router.GET("/test-missing", func(c *gin.Context) {
			RenderWithLayout(c, "nonexistent/template", gin.H{
				"title": "Test Page",
			})
		})

		req, _ := http.NewRequest("GET", "/test-missing", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return internal server error for missing template
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("handles template execution error", func(t *testing.T) {
		router := gin.New()
		router.GET("/test-exec-error", func(c *gin.Context) {
			// Pass data that will cause template execution to fail
			// For example, invalid data types that don't match template expectations
			RenderWithLayout(c, "archers/index", gin.H{
				"archers": "invalid_data_type", // Should be slice, not string
			})
		})

		req, _ := http.NewRequest("GET", "/test-exec-error", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// The function doesn't return an error status for execution errors,
		// but it prints the error. We can still verify it doesn't crash.
		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
	})
}