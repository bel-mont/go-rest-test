package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/infrastructure/auth"
	"net/http"
)

// AuthMiddleware checks for the presence of a valid JWT token in cookies and sets the user ID in the request context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the cookie
		token, err := c.Cookie(auth.AuthTokenCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Verify the token
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store the user ID in context for the next handler
		c.Set(auth.UserIDContextKey, claims.UserID)
		c.Next()
	}
}
