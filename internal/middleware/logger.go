package middleware
import (
	"github.com/gin-gonic/gin"
	"time"
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

		c.Next()
	}
}
