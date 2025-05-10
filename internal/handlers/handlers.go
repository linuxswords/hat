package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":   "HAT - Home",
		"Content": "Welcome to the Archery Tournament System!",
	})
}

func ShowParticipantsPage(c *gin.Context) {
	// Render a page to add participants
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":   "HAT - Participants",
		"Content": "Participants Page",
	})
}

func AddParticipant(c *gin.Context) {
	// Logic to add a participant
	c.String(http.StatusOK, "Participant added")
}

func ShowScoresPage(c *gin.Context) {
	// Render a page to add scores
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":   "HAT - Scores",
		"Content": "Scores Page",
	})
}

func AddScore(c *gin.Context) {
	// Logic to add a score
	c.String(http.StatusOK, "Score added")
}

func ShowResultsPage(c *gin.Context) {
	// Render a page to show results
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":   "HAT - Results",
		"Content": "Results Page",
	})
}

func ShowHandicapsPage(c *gin.Context) {
	// Render a page to manage handicaps
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":   "HAT - Handicaps",
		"Content": "Handicaps Page",
	})
}

func AddHandicap(c *gin.Context) {
	// Logic to add a handicap
	c.String(http.StatusOK, "Handicap added")
}
