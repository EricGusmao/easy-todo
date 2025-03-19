package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/EricGusmao/easy-todo/auth"
	"github.com/EricGusmao/easy-todo/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	userRepo := user.NewRepository(db)
	authService := auth.NewService(userRepo)

	http.HandleFunc("POST /signup", auth.HandleSignup(authService))

	http.ListenAndServe(":8080", nil)
}
