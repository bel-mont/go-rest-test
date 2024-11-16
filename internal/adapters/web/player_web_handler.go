package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/auth"
	"go-rest-test/pkg/html"
	"log"
	"net/http"
)

type PlayerWebHandler struct {
	repo repository.PlayerRepository
}

func NewPlayerWebHandler(repo repository.PlayerRepository) *PlayerWebHandler {
	return &PlayerWebHandler{repo: repo}
}

// RenderPlayersList renders the player list as an HTML page.
func (h *PlayerWebHandler) RenderPlayersList(c *gin.Context) {
	// Parse the players list template with header and footer
	tmpl, err := html.BaseLayoutTemplate("web/views/players/list.gohtml")
	if err != nil {
		log.Printf("Error loading players list template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	// Render the signup form template
	isAuthenticated := auth.IsUserAuthenticated(c)

	// Retrieve players from the repository
	players, err := h.repo.GetAllPlayers(context.Background())
	if err != nil {
		err := tmpl.ExecuteTemplate(c.Writer, "players/list.gohtml", gin.H{"error": "Failed to retrieve players"})
		if err != nil {
			log.Printf("Error rendering players list template: %v", err)
			return
		}
		return
	}

	data := gin.H{
		"players":           players,
		"title":             "Players",
		"header":            "Players",
		"UserAuthenticated": isAuthenticated,
	}
	// Render the template with players data
	err = tmpl.ExecuteTemplate(c.Writer, "players/list.gohtml", data)
	if err != nil {
		log.Printf("Error rendering players list template: %v", err)
		return
	}
}
