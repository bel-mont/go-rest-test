package db

import "time"

type Player struct {
	ID           int
	Username     string
	Level        int
	TotalMatches int
	TotalWins    int
	LastLogin    time.Time
}
