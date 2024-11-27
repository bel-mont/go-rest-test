package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"net/http"
)

type PlayerHandler struct {
	repo repository.Repository[entities.Player]
}

func NewPlayerHandler(repo repository.Repository[entities.Player]) PlayerHandler {
	return PlayerHandler{repo: repo}
}

func (h PlayerHandler) CreatePlayer(c *gin.Context) {
	var player entities.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.repo.Create(context.Background(), player)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h PlayerHandler) GetPlayerByID(c *gin.Context) {
	id := c.Param("id")
	player, err := h.repo.Get(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (h PlayerHandler) UpdatePlayer(c *gin.Context) {
	var player entities.Player

	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.Update(context.Background(), player)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player updated successfully"})
}

func (h PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}

func (h PlayerHandler) GetAllPlayers(c *gin.Context) {
	players, err := h.repo.Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve players"})
		return
	}

	c.JSON(http.StatusOK, players)
}
