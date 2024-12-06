package main

import (
	"main"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default();

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200,"Hello World!");
	})

	router.Run();
}
