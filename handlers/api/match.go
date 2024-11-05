package api

import "github.com/gin-gonic/gin"

// SubmitMatch handler
func SubmitMatch(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Submit matchView endpoint"})
}

// Matchmaking handler
func Matchmaking(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Matchmaking endpoint"})
}
