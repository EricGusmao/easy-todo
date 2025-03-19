package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetById(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &userRepositoryImpl{
		db: db,
	}
}

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func (u *userRepositoryImpl) Create(ctx context.Context, user *User) (*User, error) {
	query := u.db.QueryRow(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", user.Email, user.PasswordHash)
	var id uint64
	err := query.Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (u *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	panic("unimplemented")
}

func (u *userRepositoryImpl) GetById(ctx context.Context, id uint64) (*User, error) {
	panic("unimplemented")
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &userRepositoryImpl{
		db: db,
	}
}
