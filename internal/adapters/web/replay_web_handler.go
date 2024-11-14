package web

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/pkg/html"
	"log"
	"net/http"
)

type ReplayWebHandler struct{}

func NewReplayWebHandler() *ReplayWebHandler {
	return &ReplayWebHandler{}
}

func (h *ReplayWebHandler) RenderIndex(c *gin.Context) {
	// Parse the signup form template with header and footer
	tmpl, err := html.LoggedLayoutTemplate("web/views/replay/index.gohtml")
	if err != nil {
		log.Printf("Error loading index replay template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	// Render the signup form template
	data := gin.H{
		"title":  "Replay List",
		"header": "Replay List",
	}
	err = tmpl.ExecuteTemplate(c.Writer, "replay/index.gohtml", data)
	if err != nil {
		log.Printf("Error executing index replay template: %v", err)
		return
	}
}
