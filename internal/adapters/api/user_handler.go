package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/auth"
	"net/http"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) UserHandler {
	return UserHandler{userRepo: userRepo}
}

// Signup handles user registration
func (h UserHandler) Signup(c *gin.Context) {
	// Define a struct for binding both JSON and form data
	var input struct {
		Email    string `form:"email" json:"email" binding:"required,email"`
		Password string `form:"password" json:"password" binding:"required,min=8"`
	}

	// Bind request data to input struct (supports JSON and form data)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if user already exists
	_, err := h.userRepo.GetUserByEmail(context.Background(), input.Email)
	if !errors.Is(err, repository.ErrItemNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user entity
	user := entities.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Save the user in the database
	newUser, err := h.userRepo.Create(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "message": err.Error()})
		return
	}

	auth.SetTokenCookies(c, newUser.ID)

	c.Header("HX-Redirect", "/replay")
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h UserHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `form:"email" json:"email" binding:"required,email"`
		Password string `form:"password" json:"password" binding:"required,min=8"`
	}

	// Bind input using `ShouldBind` to support both JSON and form data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userRepo.GetUserByEmail(context.Background(), input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}

	matches, err := auth.CheckPasswordHash(input.Password, user.Password)
	if err != nil || !matches {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password"})
		return
	}

	auth.SetTokenCookies(c, user.ID)

	// Set HX-Redirect header to trigger a redirect in HTMX
	c.Header("HX-Redirect", "/replay")
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h UserHandler) Logout(c *gin.Context) {
	// Set the "auth_token" cookie to expire immediately
	c.SetCookie(auth.AuthTokenCookieName, "", -1, "/", "", true, true) // Negative max-age to delete the cookie

	// Check if the request is from HTMX
	if c.GetHeader("HX-Request") != "" {
		c.Header("HX-Redirect", "/")
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
