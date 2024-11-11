package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/repository"
	"html/template"
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
	tmpl, err := template.ParseFiles(
		"web/views/players/list.gohtml",
		"web/views/layouts/base-header.gohtml",
		"web/views/layouts/base-footer.gohtml",
	)
	if err != nil {
		log.Printf("Error loading players list template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

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
		"players": players,
		"title":   "Players",
		"header":  "Players",
	}
	// Render the template with players data
	err = tmpl.ExecuteTemplate(c.Writer, "players/list.gohtml", data)
	if err != nil {
		log.Printf("Error rendering players list template: %v", err)
		return
	}
}
