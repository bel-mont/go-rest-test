package db

import "time"

// Auth represents a row in the auth table
type Auth struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
