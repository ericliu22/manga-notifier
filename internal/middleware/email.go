package middleware

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"
)

func EmailClient() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mailClient, err := setupMailClient()
		if err != nil {
			log.Fatalf("Failed to setup mailClient: %s", err.Error())
		}

		c := context.WithValue(ctx.Request.Context(), "mailClient", mailClient)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

func setupMailClient() (*mail.Client, error) {
	username := os.Getenv("EMAIL_USERNAME")
	if (username == "") {
		return nil, errors.New("Failed to get env var EMAIL_USERNAME")
	}

	password := os.Getenv("EMAIL_PASSWORD")
	if (password == "") {
		return nil, errors.New("Failed to get env var EMAIL_PASSWORD")
	}

	client, err := mail.NewClient(
		"manganotifier", //This will be our domain like manganotifier.com
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
