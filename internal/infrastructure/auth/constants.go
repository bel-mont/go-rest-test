package auth

const (

	// AuthTokenCookieName holds the name of the cookie used for authentication tokens.
	AuthTokenCookieName = "auth_token"

	// UserAuthenticatedKey is used as the context key for storing the authentication status of a user.
	UserAuthenticatedKey = "UserAuthenticated"

	// UserIDContextKey is the context key used to store the user's unique identifier.
	UserIDContextKey = "userID"
)
