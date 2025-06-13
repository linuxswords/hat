package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/handlers"
	"github.com/linuxswords/hat/internal/models/bootstrap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("hat.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := bootstrap.BootstrapData(db); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.StaticFile("favicon.ico", "./static/images/favicon.ico")
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")
	// Define routes
	r.GET("/", handlers.ShowHomePage)
	r.GET("/archers", handlers.ShowArchersPage)
	r.GET("/archers/edit/:id", handlers.ShowArchersPage)
	r.GET("/scores", handlers.ShowScoresPage)
	r.GET("/resultLists", handlers.ShowResultLists)

	// BowClass routes
	r.GET("/bowclasses", handlers.ShowBowClassesPage)
	r.PUT("/bowclasses/:id", handlers.UpdateBowClass)

	// Handicap routes
	r.GET("/handicaps", handlers.ShowHandicapsPage)
	r.GET("/handicaps/add", handlers.AddHandicapSet)
	r.POST("/handicaps/add", handlers.AddHandicapSet)
	r.GET("/handicaps/edit/:id", handlers.EditHandicapSet)
	r.POST("/handicaps/edit/:id", handlers.EditHandicapSet)
	r.DELETE("/handicaps/:id", handlers.DeleteHandicap)

	// Tournament routes
	r.GET("/tournaments", handlers.ShowTournamentsPage)
	r.POST("/tournaments", handlers.AddTournament)        // For adding new tournaments
	r.POST("/tournaments/:id", handlers.UpdateTournament) // For updating existing tournaments
	r.DELETE("/tournaments/:id", handlers.DeleteTournament)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/archers/tournament/:id", handlers.GetArchers)
		api.GET("/archers/:id", handlers.GetArcher)
		api.POST("/archers", handlers.AddArcher)
		api.DELETE("/archers/:id", handlers.DeleteArcher)

		api.GET("/tournaments/:id/download", handlers.DownloadTournamentPDF)
		api.GET("/tournaments/:id", handlers.GetTournament)

		api.POST("/scores", handlers.SaveScores)
		api.DELETE("/bowclasses/:id", handlers.DeleteBowClass)
		api.POST("/bowclasses", handlers.AddBowClass)
	}

	r.Run(":8080")
}
