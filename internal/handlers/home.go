package handlers

import (
	"github.com/gin-gonic/gin"
)

func HomeHandler(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.String(400, "Not found")
		return
	}

	ctx.String(200, "Hello world")

}
