package main

import (
	"manga-notifier/internal/routes"
	"manga-notifier/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default();

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200,"Hello World!");
	})

	routes.SetupCoreRoutes(asdf)

	router.Run();
}
