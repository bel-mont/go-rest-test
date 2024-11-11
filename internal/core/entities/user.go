package entities

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't expose password in JSON responses
	CreatedAt time.Time `json:"created_at"`
}
