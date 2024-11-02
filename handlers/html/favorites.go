package html

import "github.com/gin-gonic/gin"

func FavoritesPage(c *gin.Context) {
	c.HTML(200, "favorites", gin.H{
		"header":      "Favorites",
		"title":       "Favorites",
		"description": "Your favorite matches!",
		"keywords":    "sf6, fighting games, favorites",
	})
}
