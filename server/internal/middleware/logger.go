package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := ctx.Writer.Status()
		log.Println(status)

		ctx.Next()
	}
}
