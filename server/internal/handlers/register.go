package handlers

import (
	"server/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func createUser(email string) *models.User {
	var user models.User

	user = models.User {
		ID: uuid.New().String(),
		Email: email,
		CreatedAt: time.Now(),
	}

	return &user
}

func RegisterUser(ctx *gin.Context) {
	body := ctx.Request.Body

}
