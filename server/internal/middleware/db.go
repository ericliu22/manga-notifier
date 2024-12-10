package middleware

import (
	"context"
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

func setupDatabase() (*pgxpool.Pool, error){
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
