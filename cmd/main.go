package main

import (
	"context"
	"os"
	"time"

	"github.com/Mangrover007/banking-backend/internals/repository"
	"github.com/Mangrover007/banking-backend/routers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var DB_URI = os.Getenv("DATABASE_URI")
	if DB_URI == "" {
		panic("no DB URI in env variables")
	}

	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	conn, _ := pgx.Connect(ctx, DB_URI)
	query := repository.New(conn)

	defer conn.Close(ctx)
	defer cancelFunc()

	routers.AuthRouter(server, query)
	routers.BankRouter(server, query, conn)

	server.Run(":8080")
}
