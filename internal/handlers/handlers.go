package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home/index", gin.H{
		"Title":   "HAT - Home",
		"Content": "Welcome to the Handicap Archery Tournament System!",
	})
}
