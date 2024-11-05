package entities

import "time"

type Player struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Level        int       `json:"level"`
	TotalMatches int       `json:"total_matches"`
	TotalWins    int       `json:"total_wins"`
	LastLogin    time.Time `json:"last_login"`
}
