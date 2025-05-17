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
	db.AutoMigrate(&models.BowClass{}, &models.Archer{}, &models.HandycapSet{}, &models.HandycapEntry{}, &models.Tournament{})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	// Define routes
	r.GET("/", handlers.ShowHomePage)
	r.GET("/archers", handlers.ShowArchersPage)
	r.POST("/archers", handlers.AddArcher)
	r.POST("/archers/edit/:id", handlers.UpdateArcher)
	r.POST("/archers/delete/:id", handlers.DeleteArcher)
	r.GET("/scores", handlers.ShowScoresPage)

	// BowClass routes
	r.GET("/bowclasses", handlers.ShowBowClassesPage)
	r.POST("/bowclasses", handlers.AddBowClass)
	r.PUT("/bowclasses/:id", handlers.UpdateBowClass)
	r.DELETE("/bowclasses/:id", handlers.DeleteBowClass)

	// Handycap routes
	r.GET("/handycaps", handlers.ShowHandycapsPage)
	r.GET("/handycaps/add", handlers.AddHandycapSet)
	r.POST("/handycaps/add", handlers.AddHandycapSet)
	r.PUT("/handycaps/:id", handlers.UpdateHandycap)
	r.DELETE("/handycaps/:id", handlers.DeleteHandycap)

	// Tournament routes
	r.GET("/tournaments", handlers.ShowTournamentsPage)
	r.POST("/tournaments", handlers.AddTournament)
	r.PUT("/tournaments/:id", handlers.UpdateTournament)
	r.DELETE("/tournaments/:id", handlers.DeleteTournament)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/archers/tournament/:id", handlers.GetArchers)
	}

	r.Run(":8080")
}
