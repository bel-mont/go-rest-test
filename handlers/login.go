package handlers

import "github.com/gin-gonic/gin"

// Login handler
func Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login endpoint"})
}
