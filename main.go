package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"

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

	http.HandleFunc("POST /signup", handleSignup)
}

type Middleware func(http.Handler) http.Handler

type UserRepository interface {
	GetById(ctx context.Context, id uint64) (*User, error)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type User struct {
	ID           uint64
	email        string
	passwordHash string
}

type CreateUserRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type UserCredentials struct {
	email    string
	password string
}

type AuthService interface {
	Signup(ctx context.Context, user *CreateUserRequest)
	Login(ctx context.Context, cred *UserCredentials)
}

type UserRepository interface {
	
}
