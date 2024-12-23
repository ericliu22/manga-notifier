package handlers

import (
	"context"
	"log"
	"net/http"
	"server/internal/middleware"
	"server/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func createUser(ctx *gin.Context, email string) error {
	var user models.User

	db, err := middleware.GetDatabasePool(ctx)
	if err != nil {
		log.Printf("Failed to get db pool: %s", err.Error())
		return err
	}

	user = models.User{
		ID:        uuid.New().String(),
		Email:     email,
		CreatedAt: time.Now(),
	}


	query := `
	INSERT INTO users (id, email, created_at)
	VALUES ($1, $2, NOW())
	RETURNING id, email, created_at
	`

	scanErr := db.QueryRow(context.Background(), query, user.ID, user.Email).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if scanErr != nil {
		log.Printf("DB Scan failed: %s", scanErr.Error())
		return scanErr
	}

	return nil
}

func RegisterUser(ctx *gin.Context) {
	var req RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Return error if JSON binding fails
		log.Printf("Bind JSON failed: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := createUser(ctx, req.Email); err != nil {
		log.Printf("Register User failed: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
