package routes

import (
	"server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCoreRoutes(router *gin.Engine) {
	router.GET("/", handlers.HomeHandler)
	SetupUserRoutes(router)
}
