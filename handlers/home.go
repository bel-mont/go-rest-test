package handlers

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.HTML(200, "layout.html", gin.H{
		"title":       "Home Page",
		"description": "This is the home page of the website.",
		"keywords":    "home, welcome, example",
		"header":      "Welcome to My Website",
		"content":     "home.html",
	})
}
