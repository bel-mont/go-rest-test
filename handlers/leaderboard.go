package handlers

import "github.com/gin-gonic/gin"

// GetLeaderboard handler
func GetLeaderboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Leaderboard endpoint"})
}
