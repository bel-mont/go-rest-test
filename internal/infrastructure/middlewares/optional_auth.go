package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/infrastructure/auth"
)

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(auth.AuthTokenCookieName)
		if err != nil || token == "" {
			// User is not authenticated
			c.Set(auth.UserAuthenticatedKey, false)
			c.Next()
			return
		}

		userID, err := auth.ValidateJWT(token)
		if err != nil {
			// Invalid token; mark as not authenticated
			c.Set(auth.UserAuthenticatedKey, false)
		} else {
			// Valid token; user is authenticated
			c.Set(auth.UserIDContextKey, userID)
			c.Set(auth.UserAuthenticatedKey, true)
		}

		c.Next()
	}
}
