package web

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/pkg/html"
	"log"
	"net/http"
)

type HomeWebHandler struct{}

func NewHomeWebHandler() *HomeWebHandler {
	return &HomeWebHandler{}
}

func (h *HomeWebHandler) RenderHome(c *gin.Context) {
	// Parse the signup form template with header and footer
	tmpl, err := html.BaseLayoutTemplate("web/views/index.gohtml")
	if err != nil {
		log.Printf("Error loading home template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	// Render the signup form template
	data := gin.H{
		"title":  "FG Replay Analyzer",
		"header": "FG Replay Analyzer",
	}
	err = tmpl.ExecuteTemplate(c.Writer, "home.gohtml", data)
	if err != nil {
		log.Printf("Error executing signup template: %v", err)
		return
	}
}
