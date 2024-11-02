package html

import "github.com/gin-gonic/gin"

func LeaderboardPage(c *gin.Context) {
	c.HTML(200, "leaderboard", gin.H{
		"header":      "Leaderboard",
		"title":       "Leaderboard",
		"description": "The leaderboard!",
		"keywords":    "sf6, fighting games, leaderboard",
	})
}
