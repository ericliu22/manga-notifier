package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetupMiddleware(router *gin.Engine) {
	router.Use(Logger())
	router.Use(EmailClient())
	router.Use(DatabaseMiddleware())
}
