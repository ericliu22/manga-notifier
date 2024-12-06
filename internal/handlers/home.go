package handlers

import (
	"github.com/gin-gonic/gin"
)

func HomeHandler(ctx *gin.Context) {
		ctx.String(200,"Hello World!");
}
