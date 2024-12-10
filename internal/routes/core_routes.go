package routes

import (
	"github.com/gin-gonic/gin"
	"manga-notifier/internal/handlers"
)

func SetupCoreRoutes(router *gin.Engine) {
	router.GET("/", handlers.HomeHandler)
}
