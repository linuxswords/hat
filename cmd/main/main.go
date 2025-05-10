package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/", showHomePage)
	r.GET("/participants", showParticipantsPage)
	r.POST("/participants", addParticipant)
	r.GET("/scores", showScoresPage)
	r.POST("/scores", addScore)
	r.GET("/results", showResultsPage)
	r.GET("/admin/handicaps", showHandicapsPage)
	r.POST("/admin/handicaps", addHandicap)

	// Start the server
	r.Run(":8080")
}
