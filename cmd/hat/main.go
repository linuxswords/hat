package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/handlers"
)

func main() {
	// Create Gin router
	r := gin.Default()

	// Serve static files (CSS, JS, images)
	r.Static("/static", "./static")

	// Routes
	r.GET("/", func(c *gin.Context) {
		handlers.RenderWithLayout(c, "home/index", gin.H{
			"title": "HAT - Handicap Archery Tournament",
		})
	})

	// Bow Classes routes
	bowclasses := r.Group("/bowclasses")
	{
		bowclasses.GET("/", handlers.BowClassesList)
		bowclasses.GET("/new", handlers.BowClassesNew)
		bowclasses.GET("/:id", handlers.BowClassesShow)
	}

	// Archers routes
	archers := r.Group("/archers")
	{
		archers.GET("/", handlers.ArchersList)
		archers.GET("/new", handlers.ArchersNew)
		archers.POST("/", handlers.ArchersCreate)
		archers.GET("/:id", handlers.ArchersShow)
		archers.GET("/:id/edit", handlers.ArchersEdit)
		archers.POST("/:id", handlers.ArchersUpdate)
		archers.POST("/:id/delete", handlers.ArchersDelete)
	}

	// Tournaments routes
	tournaments := r.Group("/tournaments")
	{
		tournaments.GET("/", handlers.TournamentsList)
		tournaments.GET("/new", handlers.TournamentsNew)
		tournaments.POST("/", handlers.TournamentsCreate)
		tournaments.GET("/:id", handlers.TournamentsShow)
		tournaments.GET("/:id/edit", handlers.TournamentsEdit)
		tournaments.POST("/:id", handlers.TournamentsUpdate)
		tournaments.POST("/:id/delete", handlers.TournamentsDelete)
	}

	// Handicaps routes
	handicaps := r.Group("/handicaps")
	{
		handicaps.GET("/", handlers.HandicapsList)
		handicaps.GET("/:id", handlers.HandicapsShow)
	}

	// Start server on port 8080
	r.Run(":8080")
}
