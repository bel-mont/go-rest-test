package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SetTokenCookies generates and sets the access and refresh tokens as secure, HTTP-only cookies
func SetTokenCookies(c *gin.Context, userID string) {
	// Generate Access Token (1 hour expiration)
	accessToken, err := GenerateToken(userID, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate Refresh Token (1 month expiration)
	const refreshTokenTTLValue = 30 * 24 * time.Hour
	refreshToken, err := GenerateToken(userID, refreshTokenTTLValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Set the tokens as secure, HTTP-only cookies
	const accessTokenTTL = 3600            // 1 hour in seconds
	const refreshTokenTTL = 30 * 24 * 3600 // 1 month in seconds

	c.SetCookie(AuthTokenCookieName, accessToken, accessTokenTTL, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, refreshTokenTTL, "/", "", true, true)
}
