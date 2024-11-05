package html

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-test/repository"
	"net/http"
)

// MatchHandler defines a handler with access to repositories
type MatchHandler struct {
	repo repository.MatchRepository
}

// NewMatchHandler creates a new MatchHandler instance
func NewMatchHandler(repo repository.MatchRepository) *MatchHandler {
	return &MatchHandler{repo: repo}
}

// MatchListPage Displays a list of matches
func (h *MatchHandler) MatchListPage(c *gin.Context) {
	matches, err := h.repo.GetMatchList(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve matches"})
		return
	}

	// Render the template with matches data
	c.HTML(http.StatusOK, "matches", gin.H{
		"title":       "SF6 Rankings",
		"description": "Welcome to SF6 Rankings, where you can follow matches, leaderboards, and participate in the community.",
		"keywords":    "sf6, rankings, matches, leaderboards",
		"header":      "SF6 Rankings",
		"Matches":     matches,
	})
}

func (h *MatchHandler) SubmitMatchPage(c *gin.Context) {
	c.HTML(200, "submit-matchView", gin.H{
		"header":      "Submit Match",
		"title":       "Submit Match",
		"description": "Submit a matchView to the SF6 Rankings database!",
		"keywords":    "sf6, fighting games, submit matchView",
	})
}

func (h *MatchHandler) MatchDetailPage(c *gin.Context) {
	c.HTML(200, "matchView-detail", gin.H{
		"header":      "Match Detail",
		"title":       "Match Detail",
		"description": "View the details of a matchView!",
		"keywords":    "sf6, fighting games, matchView detail",
	})
}

func InitRoutes(router *gin.Engine, matchHandler *MatchHandler) {
	router.GET("/matches", matchHandler.MatchListPage)
}
