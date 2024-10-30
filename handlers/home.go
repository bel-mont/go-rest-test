package handlers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.HTML(200, "layout.html", gin.H{
		"title":       "SF6 Rankings",
		"description": "Welcome to SF6 Rankings, where you can follow matches, leaderboards, and participate in the community.",
		"keywords":    "sf6, rankings, matches, leaderboards",
		"header":      "SF6 Rankings",
		"content":     "home.html",
	})
}
