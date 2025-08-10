package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/models"
	"github.com/linuxswords/hat/internal/repositories"
)

// ArcherRepository defines the contract for archer data operations
type ArcherRepository interface {
	GetAll() []models.Archer
	GetByID(id int) (*models.Archer, error)
	Create(archer models.Archer) (*models.Archer, error)
	Update(id int, archer models.Archer) (*models.Archer, error)
	Delete(id int) error
	GetByBowClass(bowClass string) []models.Archer
	GetByGender(gender string) []models.Archer
}

var archerRepo ArcherRepository = repositories.NewArcherRepository()

// loadBowClassesForArchers loads bow classes for dropdowns in archer forms
func loadBowClassesForArchers() ([]BowClass, error) {
	return loadBowClassesFromJSON()
}

// ArchersList handles GET /archers
func ArchersList(c *gin.Context) {
	archers := archerRepo.GetAll()
	RenderWithLayout(c, "archers/index", gin.H{
		"title":   "Archers",
		"archers": archers,
	})
}

// ArchersNew handles GET /archers/new
func ArchersNew(c *gin.Context) {
	bowClasses, err := loadBowClassesForArchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load bow classes"})
		return
	}

	RenderWithLayout(c, "archers/new", gin.H{
		"title":      "Create New Archer",
		"bowClasses": bowClasses,
	})
}

// ArchersCreate handles POST /archers
func ArchersCreate(c *gin.Context) {
	name := c.PostForm("name")
	gender := c.PostForm("gender")
	bowClass := c.PostForm("bow_class")
	email := c.PostForm("email")

	// Basic validation
	if name == "" || gender == "" || bowClass == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Create new archer
	archer := models.Archer{
		Name:     name,
		Gender:   gender,
		BowClass: bowClass,
		Email:    email,
	}

	// Add to repository
	_, err := archerRepo.Create(archer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create archer"})
		return
	}

	// Redirect to archer list
	c.Redirect(http.StatusFound, "/archers")
}

// ArchersShow handles GET /archers/:id
func ArchersShow(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid archer ID"})
		return
	}

	// Find archer by ID
	archer, err := archerRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}

	RenderWithLayout(c, "archers/show", gin.H{
		"title":  "Archer Details",
		"archer": archer,
	})
}

// ArchersEdit handles GET /archers/:id/edit
func ArchersEdit(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid archer ID"})
		return
	}

	// Find archer by ID
	archer, err := archerRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}

	bowClasses, err := loadBowClassesForArchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load bow classes"})
		return
	}

	RenderWithLayout(c, "archers/edit", gin.H{
		"title":      "Edit Archer",
		"archer":     archer,
		"bowClasses": bowClasses,
	})
}

// ArchersUpdate handles POST /archers/:id (with _method=PUT)
func ArchersUpdate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid archer ID"})
		return
	}

	name := c.PostForm("name")
	gender := c.PostForm("gender")
	bowClass := c.PostForm("bow_class")
	email := c.PostForm("email")

	// Basic validation
	if name == "" || gender == "" || bowClass == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Update archer
	updatedArcher := models.Archer{
		Name:     name,
		Gender:   gender,
		BowClass: bowClass,
		Email:    email,
	}

	_, err = archerRepo.Update(id, updatedArcher)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}

	// Redirect to archer details
	c.Redirect(http.StatusFound, "/archers/"+idParam)
}

// ArchersDelete handles POST /archers/:id/delete
func ArchersDelete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid archer ID"})
		return
	}

	// Delete archer from repository
	err = archerRepo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archer not found"})
		return
	}

	// Redirect to archer list
	c.Redirect(http.StatusFound, "/archers")
}

