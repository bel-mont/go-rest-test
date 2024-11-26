package entities

import "time"

// User represents an application user with ID, Email, Password, and CreatedAt fields.
type User struct {
	// ID represents the unique identifier for the user.
	ID string `dynamodbav:"id" json:"id"`
	// Email represents the user's email address and must be unique.
	Email string `dynamodbav:"email" json:"email"`
	// Password stores the user's password, excluded from JSON serialization for security reasons.
	Password string `dynamodbav:"password" json:"-"`
	// CreatedAt represents the timestamp when the user was created, stored in DynamoDB as 'created_at' and JSON as 'created_at'.
	CreatedAt time.Time `dynamodbav:"created_at" json:"created_at"`
}

func (u User) GetID() string {
	return u.ID
}

func (u User) SetID(id string) {
	u.ID = id
}
