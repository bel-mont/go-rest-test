package html

import (
	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	c.HTML(200, "login", gin.H{
		"header":      "SF6 Rankings",
		"title":       "SF6 Rankings - Login",
		"description": "SF6 rankings, matches, and leaderboards!",
		"keywords":    "sf6, fighting games, leaderboards",
	})
}
