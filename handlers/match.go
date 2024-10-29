package handlers

import (
	"github.com/gin-gonic/gin"
)

// SubmitMatch handler
func SubmitMatch(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Submit match endpoint"})
}

// Matchmaking handler
func Matchmaking(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Matchmaking endpoint"})
}
