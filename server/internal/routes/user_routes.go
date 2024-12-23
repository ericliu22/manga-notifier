package routes

import (
	"server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	router.POST("/register", handlers.RegisterUser)
}
