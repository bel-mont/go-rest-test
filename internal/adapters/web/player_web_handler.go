package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/repository"
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
	players, err := h.repo.GetAllPlayers(context.Background())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "players/list.gohtml", gin.H{"error": "Failed to retrieve players"})
		return
	}

	c.HTML(http.StatusOK, "players/list.gohtml", players)
}
