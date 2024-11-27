package web

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/entities"
	repositoryInterface "go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/auth"
	"go-rest-test/pkg/html"
	"log"
	"net/http"
)

type ReplayWebHandler struct {
	replayRepo repositoryInterface.Repository[entities.Replay]
}

func NewReplayWebHandler(replayRepo repositoryInterface.Repository[entities.Replay]) ReplayWebHandler {
	return ReplayWebHandler{
		replayRepo: replayRepo,
	}
}

func (h ReplayWebHandler) RenderIndex(c *gin.Context) {
	tmpl, err := html.BaseLayoutTemplate("web/views/replay/index.gohtml")
	if err != nil {
		log.Printf("Error loading index replay template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	// Fetch all replays
	replays, err := h.replayRepo.Scan(c.Request.Context())
	if err != nil {
		log.Printf("Error fetching replays: %v", err)
		replays = []entities.Replay{} // Use empty slice if error
	}

	isAuthenticated := auth.IsUserAuthenticated(c)

	data := gin.H{
		"title":             "Replay List",
		"header":            "Replay List",
		"UserAuthenticated": isAuthenticated,
		"Replays":           replays,
	}

	err = tmpl.ExecuteTemplate(c.Writer, "replay/index.gohtml", data)
	if err != nil {
		log.Printf("Error executing index replay template: %v", err)
		return
	}
}

func (h ReplayWebHandler) RenderUploadPage(c *gin.Context) {
	tmpl, err := html.BaseLayoutTemplate("web/views/replay/upload.gohtml")
	if err != nil {
		log.Printf("Error loading upload replay template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	isAuthenticated := auth.IsUserAuthenticated(c)

	data := gin.H{
		"title":             "Upload Replay",
		"header":            "Upload Replay",
		"UserAuthenticated": isAuthenticated,
	}

	err = tmpl.ExecuteTemplate(c.Writer, "replay/upload.gohtml", data)
	if err != nil {
		log.Printf("Error executing upload replay template: %v", err)
		return
	}
}
