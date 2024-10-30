package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
)

// GetPlayerStats handler
func GetPlayerStats(c *gin.Context) {
	username := c.Param("id")

	var player struct {
		Username     string `json:"username"`
		Level        int    `json:"level"`
		TotalMatches int    `json:"total_matches"`
		TotalWins    int    `json:"total_wins"`
	}

	// Query for a player's stats based on their username
	err := db.QueryRow(context.Background(), `
		SELECT username, level, total_matches, total_wins
		FROM players
		WHERE username = $1
	`, username).Scan(&player.Username, &player.Level, &player.TotalMatches, &player.TotalWins)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve player data"})
		}
		return
	}

	c.JSON(http.StatusOK, player)
}

// GetPlayers handler retrieves all players from the database
func GetPlayers(c *gin.Context) {
	rows, err := db.Query(context.Background(), `
		SELECT username, level, total_matches, total_wins
		FROM players
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve players"})
		return
	}
	defer rows.Close()

	var players []struct {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse player data"})
			return
		}
		players = append(players, player)
	}

	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading player data"})
		return
	}

	c.JSON(http.StatusOK, players)
}
