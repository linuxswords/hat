package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxswords/hat/internal/handlers"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	// Define routes
	r.GET("/", handlers.ShowHomePage)
	r.GET("/participants", handlers.ShowParticipantsPage)
	r.POST("/participants", handlers.AddParticipant)
	r.GET("/scores", handlers.ShowScoresPage)
	r.POST("/scores", handlers.AddScore)
	r.GET("/results", handlers.ShowResultsPage)
	r.GET("/admin/handicaps", handlers.ShowHandicapsPage)
	r.POST("/admin/handicaps", handlers.AddHandicap)

	// Start the server
	r.Run(":8987")
}
