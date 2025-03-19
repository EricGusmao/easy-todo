package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/EricGusmao/easy-todo/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Signup(ctx context.Context, r *CreateUserRequest) (string, error)
	Login(ctx context.Context, r *LoginUserRequest) (string, error)
	UserFromToken(ctx context.Context, token string) (*user.User, error)
}

type service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) Signup(ctx context.Context, r *CreateUserRequest) (string, error) {
	if r.Password == "" {
		return "", errors.New("password is required")
	}

	if r.Email == "" {
		return "", errors.New("email is required")
	}

	if r.Password != r.PasswordConfirmation {
		return "", errors.New("passwords do not match")
	}

	hashedPassword, err := hashPassword(r.Password)
	if err != nil {
		return "", err
	}

	newUser, err := s.userRepo.Create(ctx, &user.User{Email: r.Email, PasswordHash: hashedPassword})
	if err != nil {
		return "", err
	}

	secureToken, err := generateTokenFor(newUser)
	if err != nil {
		return "", err
	}

	return secureToken, nil
}

func (s *service) Login(ctx context.Context, r *LoginUserRequest) (string, error) {
	if r.Password == "" {
		return "", errors.New("password is required")
	}

	if r.Email == "" {
		return "", errors.New("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, r.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password))
	if err != nil {
		return "", err
	}

	token, err := generateTokenFor(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) UserFromToken(ctx context.Context, tokenString string) (*user.User, error) {
	if tokenString == "" {
		return nil, errors.New("token is an empty string")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong Algorithm")
		}
		return os.Getenv("SECRET_KEY"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid JWT")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	userID, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetById(ctx, userID)
}

func generateTokenFor(newUser *user.User) (string, error) {
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  strconv.FormatUint(newUser.ID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secureToken, err := token.SignedString(os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}
	return secureToken, nil
}

func hashPassword(rawPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
