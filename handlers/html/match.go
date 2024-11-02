package html

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// MatchListPage Displays a list of matches
func MatchListPage(c *gin.Context) {
	rows, err := db.Query(context.Background(), `
		SELECT p1.username AS player1, p2.username AS player2, 
		       CASE WHEN m.winner_id = m.player1_id THEN p1.username ELSE p2.username END AS winner
		FROM matches m
		JOIN players p1 ON m.player1_id = p1.id
		JOIN players p2 ON m.player2_id = p2.id
		ORDER BY m.id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve matches"})
		return
	}
	defer rows.Close()

	// Define a struct for match data
	var matches []struct {
		Player1 string `json:"player1"`
		Player2 string `json:"player2"`
		Winner  string `json:"winner"`
	}

	// Iterate over rows and scan data
	for rows.Next() {
		var match struct {
			Player1 string `json:"player1"`
			Player2 string `json:"player2"`
			Winner  string `json:"winner"`
		}
		if err := rows.Scan(&match.Player1, &match.Player2, &match.Winner); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse match data"})
			return
		}
		matches = append(matches, match)
	}

	// Check for row errors
	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading match data"})
		return
	}

	// Render the template with matches data
	c.HTML(200, "matches", gin.H{
		"title":       "SF6 Rankings",
		"description": "Welcome to SF6 Rankings, where you can follow matches, leaderboards, and participate in the community.",
		"keywords":    "sf6, rankings, matches, leaderboards",
		"header":      "SF6 Rankings",
		"Matches":     matches,
	})
}

func SubmitMatchPage(c *gin.Context) {
	c.HTML(200, "submit-match", gin.H{
		"header":      "Submit Match",
		"title":       "Submit Match",
		"description": "Submit a match to the SF6 Rankings database!",
		"keywords":    "sf6, fighting games, submit match",
	})
}

func MatchDetailPage(c *gin.Context) {
	c.HTML(200, "match-detail", gin.H{
		"header":      "Match Detail",
		"title":       "Match Detail",
		"description": "View the details of a match!",
		"keywords":    "sf6, fighting games, match detail",
	})
}
