package routes

import (
	"github.com/gin-gonic/gin"
	"server/internal/handlers"
)

func SetupUserRoutes(router *gin.Engine) {
	router.POST("/register", handlers.RegisterUser)
}
