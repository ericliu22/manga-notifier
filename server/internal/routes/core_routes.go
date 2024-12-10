package routes

import (
	"github.com/gin-gonic/gin"
	"server/internal/handlers"
)

func SetupCoreRoutes(router *gin.Engine) {
	router.GET("/", handlers.HomeHandler)
	SetupUserRoutes(router)
}
