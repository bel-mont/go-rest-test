package web

import (
	"go-rest-test/internal/infrastructure/auth"
	"go-rest-test/pkg/html"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserWebHandler struct{}

func NewUserWebHandler() UserWebHandler {
	return UserWebHandler{}
}

// RenderSignupForm renders the signup form HTML page.
func (h UserWebHandler) RenderSignupForm(c *gin.Context) {
	// Parse the signup form template with header and footer
	tmpl, err := html.BaseLayoutTemplate("web/views/users/signup.gohtml")
	if err != nil {
		log.Printf("Error loading signup template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	// Render the signup form template
	isAuthenticated := auth.IsUserAuthenticated(c)

	// Render the signup form template
	data := gin.H{
		"title":             "Sign Up",
		"header":            "Sign Up",
		"UserAuthenticated": isAuthenticated,
	}
	err = tmpl.ExecuteTemplate(c.Writer, "users/signup.gohtml", data)
	if err != nil {
		log.Printf("Error executing signup template: %v", err)
		return
	}
}

// RenderLoginForm renders the login form HTML page.
func (h UserWebHandler) RenderLoginForm(c *gin.Context) {
	// Parse the login form template with header and footer
	tmpl, err := template.ParseFiles(
		"web/views/users/login.gohtml",
		"web/views/layouts/base-header.gohtml",
		"web/views/layouts/base-footer.gohtml",
	)
	if err != nil {
		log.Printf("Error loading login template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	data := gin.H{
		"title":  "Log In",
		"header": "Log In",
	}
	// Render the login form template
	err = tmpl.ExecuteTemplate(c.Writer, "users/login.gohtml", data)
	if err != nil {
		log.Printf("Error executing login template: %v", err)
		return
	}
}
