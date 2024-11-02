package handlers

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.HTML(200, "home", gin.H{
		"header":      "SF6 Rankings",
		"title":       "SF6 Rankings",
		"description": "SF6 rankings, matches, and leaderboards!",
		"keywords":    "sf6, fighting games, leaderboards",
	})
}
