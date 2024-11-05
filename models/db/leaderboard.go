package db

import "time"

// Leaderboard represents a row in the leaderboard table
type Leaderboard struct {
	ID        int
	PlayerID  int
	Score     int
	Rank      int
	UpdatedAt time.Time
}
