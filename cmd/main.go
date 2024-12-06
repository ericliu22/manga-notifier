package main

import (
	"manga-notifier/internal/routes"
	"manga-notifier/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default();

    middleware.SetupMiddleware(router)
	routes.SetupCoreRoutes(router)

	router.Run();
}
