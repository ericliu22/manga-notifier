package main

import (
	"manga-notifier/internal/middleware"
	"manga-notifier/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	middleware.SetupMiddleware(router)
	routes.SetupCoreRoutes(router)

	router.Run()
}
