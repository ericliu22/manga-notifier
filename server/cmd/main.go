package main

import (
	"server/internal/middleware"
	"server/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	middleware.SetupMiddleware(router)
	routes.SetupCoreRoutes(router)

	router.Run()
}
