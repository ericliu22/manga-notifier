package routes

import (
    "manga-notifier/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupCoreRoutes(router *gin.Engine) {
	router.GET("/", handlers.HomeHandler)
}
