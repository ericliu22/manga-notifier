package middleware

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)


func DatabaseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db, err := setupDatabase()
		if err != nil {
			log.Fatalf("Error setting up database: %s", err.Error())
		}

		ctx.Set("db", db)
		ctx.Next()
	}
}

func GetDatabasePool(ctx *gin.Context) (*pgxpool.Pool, error) {
	var dbpool *pgxpool.Pool

	dbpool = ctx.Value("db").(*pgxpool.Pool)
	if dbpool == nil {
		return nil, errors.New("Failed to get db from context")
	}

	return dbpool, nil
}

func setupDatabase() (*pgxpool.Pool, error){
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
