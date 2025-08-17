package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/handlers"
	"github.com/linuxswords/hat/internal/repositories"
)

func main() {
	// Create Gin router
	r := gin.Default()

	// Serve static files (CSS, JS, images)
	r.Static("/static", "./static")

	// Initialize repositories
	archerRepo := repositories.NewArcherRepository()
	tournamentRepo := repositories.NewTournamentRepository()
	handicapRepo := repositories.NewHandicapRepository()
	scoreRepo := repositories.NewScoreRepository()
	participationRepo := repositories.NewTournamentParticipationRepository()

	// Initialize handler structs with injected dependencies
	archerHandlers := &handlers.ArcherHandlers{
		ArcherRepo: archerRepo,
	}

	tournamentHandlers := &handlers.TournamentHandlers{
		TournamentRepo: tournamentRepo,
		HandicapRepo:   handicapRepo,
	}

	tournamentScoreHandlers := &handlers.TournamentScoreHandlers{
		ScoreRepo:         scoreRepo,
		TournamentRepo:    tournamentRepo,
		ArcherRepo:        archerRepo,
		HandicapRepo:      handicapRepo,
		ParticipationRepo: participationRepo,
	}

	tournamentArcherHandlers := &handlers.TournamentArcherHandlers{
		TournamentRepo:    tournamentRepo,
		ArcherRepo:        archerRepo,
		ParticipationRepo: participationRepo,
	}

	handicapHandlers := &handlers.HandicapHandlers{
		HandicapRepo: handicapRepo,
	}

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
		archers.GET("/", archerHandlers.ArchersList)
		archers.GET("/new", archerHandlers.ArchersNew)
		archers.POST("/", archerHandlers.ArchersCreate)
		archers.GET("/:id", archerHandlers.ArchersShow)
		archers.GET("/:id/edit", archerHandlers.ArchersEdit)
		archers.POST("/:id", archerHandlers.ArchersUpdate)
		archers.POST("/:id/delete", archerHandlers.ArchersDelete)
	}

	// Tournaments routes
	tournaments := r.Group("/tournaments")
	{
		tournaments.GET("/", tournamentHandlers.TournamentsList)
		tournaments.GET("/new", tournamentHandlers.TournamentsNew)
		tournaments.POST("/", tournamentHandlers.TournamentsCreate)
		tournaments.GET("/:id", tournamentHandlers.TournamentsShow)
		tournaments.GET("/:id/edit", tournamentHandlers.TournamentsEdit)
		tournaments.POST("/:id", tournamentHandlers.TournamentsUpdate)
		tournaments.POST("/:id/delete", tournamentHandlers.TournamentsDelete)
		
		// Tournament archers routes
		tournaments.GET("/:id/archers", tournamentArcherHandlers.TournamentArchersIndex)
		tournaments.GET("/:id/archers/add", tournamentArcherHandlers.TournamentArchersAdd)
		tournaments.POST("/:id/archers", tournamentArcherHandlers.TournamentArchersCreate)
		tournaments.POST("/:id/archers/:archer_id/remove", tournamentArcherHandlers.TournamentArchersRemove)
		tournaments.POST("/:id/archers/:archer_id/status", tournamentArcherHandlers.TournamentArchersUpdateStatus)
		
		// Tournament scores routes
		tournaments.GET("/:id/scores", tournamentScoreHandlers.TournamentScoresIndex)
		tournaments.GET("/:id/scores/edit", tournamentScoreHandlers.TournamentScoresEdit)
		tournaments.POST("/:id/scores", tournamentScoreHandlers.TournamentScoresUpdate)
		tournaments.GET("/:id/rankings", tournamentScoreHandlers.TournamentScoresRankings)
	}

	// Handicaps routes
	handicaps := r.Group("/handicaps")
	{
		handicaps.GET("/", handicapHandlers.HandicapsList)
		handicaps.GET("/:id", handicapHandlers.HandicapsShow)
	}

	// Start server on port 8080
	r.Run(":8080")
}
