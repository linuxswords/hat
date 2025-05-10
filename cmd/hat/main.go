package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/handlers"
	"github.com/linuxswords/hat/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("hat.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.BowClass{}, &models.Participant{}, &models.HandicapSet{}, &models.Handicap{})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.LoadHTMLGlob("templates/*")
	// Define routes
	r.GET("/", handlers.ShowHomePage)
	r.GET("/participants", handlers.ShowParticipantsPage)
	r.POST("/participants", handlers.AddParticipant)
	r.POST("/participants/edit/:id", handlers.UpdateParticipant)
	r.POST("/participants/delete/:id", handlers.DeleteParticipant)
	r.POST("/scores", handlers.AddScore)
	r.GET("/results", handlers.ShowResultsPage)

	// BowClass routes
	r.GET("/bowclasses", handlers.ShowBowClassesPage)
	r.POST("/bowclasses", handlers.AddBowClass)
	r.PUT("/bowclasses/:id", handlers.UpdateBowClass)
	r.DELETE("/bowclasses/:id", handlers.DeleteBowClass)

	r.Run(":8987")
}
