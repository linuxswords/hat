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
	r.GET("/tournament/score/:id", handlers.ShowTournamentScores)
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

	// Tournament routes
	r.GET("/tournaments", handlers.ShowTournamentsPage)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/archers/:id", handlers.GetArcher)
		api.POST("/archers", handlers.AddArcher)
		api.DELETE("/archers/:id", handlers.DeleteArcher)

		api.GET("/tournaments/:id/download", handlers.DownloadTournamentPDF)
		api.POST("/tournaments", handlers.AddTournament)
		api.GET("/tournaments/:id", handlers.GetTournament)
		api.PUT("/tournaments/:id", handlers.AddTournament)
		api.DELETE("/tournaments/:id", handlers.DeleteTournament)

		api.PUT("/scores/tournament/:tournamentId/archer/:archerID/score/:scoreID", handlers.UpdateArcherScore)
		api.DELETE("/bowclasses/:id", handlers.DeleteBowClass)
		api.POST("/bowclasses", handlers.AddBowClass)

		api.DELETE("/handicaps/:id", handlers.DeleteHandicap)
	}

	r.Run(":8080")
}
