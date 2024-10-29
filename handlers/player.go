package handlers

import "github.com/gin-gonic/gin"

// GetPlayerStats handler
func GetPlayerStats(c *gin.Context) {
	// Example of retrieving a URL parameter
	playerID := c.Param("id")
	c.JSON(200, gin.H{"message": "Player stats endpoint", "playerID": playerID})
}
