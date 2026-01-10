package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Mangrover007/banking-backend/app/internals/repository"
	"github.com/Mangrover007/banking-backend/app/routers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	var DB_HOST = os.Getenv("DB_HOST")
	var DB_USER = os.Getenv("DB_USER")
	var DB_PASSWORD = os.Getenv("DB_PASSWORD")

	if DB_HOST == "" {
		DB_HOST = "localhost:5432"
	}
	if DB_USER == "" {
		DB_USER = "postgres"
	}
	if DB_PASSWORD == "" {
		DB_PASSWORD = "pass"
	}

	var DB_URI = fmt.Sprintf("postgres://%s:%s@%s/banking", DB_USER, DB_PASSWORD, DB_HOST)

	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	conn, err := pgx.Connect(ctx, DB_URI)
	if err != nil {
		fmt.Printf("Could not connect to DB %s\n", DB_URI)
		panic(fmt.Errorf("%+w", err))
	}
	query := repository.New(conn)

	defer conn.Close(ctx)
	defer cancelFunc()

	routers.AuthRouter(server, query)
	routers.BankRouter(server, query, conn)
	routers.AccountRouter(server, conn, query)

	server.Run(":8080")
}
