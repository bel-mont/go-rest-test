package entities

import "time"

// Player represents a player in the game.
type Player struct {
	// ID is the unique identifier of the player.
	ID string `dynamodbav:"id" json:"id"`
	// Username is the name of the player.
	Username string `dynamodbav:"username" json:"username"`
	// Level indicates the player's level.
	Level int `dynamodbav:"level" json:"level"`
	// TotalMatches is the total number of matches the player has played.
	TotalMatches int `dynamodbav:"total_matches" json:"total_matches"`
	// TotalWins is the total number of matches the player has won.
	TotalWins int `dynamodbav:"total_wins" json:"total_wins"`
	// LastLogin is the last login time of the player.
	LastLogin time.Time `dynamodbav:"last_login" json:"last_login"`
}
