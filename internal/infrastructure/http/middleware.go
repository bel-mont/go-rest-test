package http

import (
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/infrastructure/middlewares"
	"log"
	"os"
	"strings"
)

func InitializeMiddlewares(router *gin.Engine) {
	router.Use(middlewares.OptionalAuthMiddleware())
	setTrustedProxies(router)
}

func setTrustedProxies(router *gin.Engine) {
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		if err := router.SetTrustedProxies(proxies); err != nil {
			log.Fatalf("Error setting trusted proxies: %v", err)
		}
	} else {
		log.Println("No trusted proxies set; all proxies are trusted by default.")
	}
}
