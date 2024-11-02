package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetLeaderboard handler
func GetLeaderboard(c *gin.Context) {
	rows, err := db.Query(context.Background(), `
		SELECT username, level, total_matches, total_wins
		FROM players
		ORDER BY total_wins DESC
		LIMIT 10
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard data"})
		return
	}
	defer rows.Close()

	var leaderboard []struct {
		Username     string `json:"username"`
		Level        int    `json:"level"`
		TotalMatches int    `json:"total_matches"`
		TotalWins    int    `json:"total_wins"`
	}

	for rows.Next() {
		var player struct {
			Username     string `json:"username"`
			Level        int    `json:"level"`
			TotalMatches int    `json:"total_matches"`
			TotalWins    int    `json:"total_wins"`
		}
		if err := rows.Scan(&player.Username, &player.Level, &player.TotalMatches, &player.TotalWins); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse leaderboard data"})
			return
		}
		leaderboard = append(leaderboard, player)
	}

	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading leaderboard data"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
