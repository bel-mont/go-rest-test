package db

import "time"

// Matchmaking represents a row in the matchmaking table
type Matchmaking struct {
	ID         int
	PlayerID   int
	SkillLevel int
	Status     string
	QueuedAt   time.Time
}
