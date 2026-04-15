package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const (
	insertUserQuery = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	selectUserByEmailQuery = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	selectUserByIDQuery = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
)

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	q := resolveQuerier(ctx, r.pool)

	if err := q.QueryRow(ctx, insertUserQuery, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt); err != nil {
		if isUniqueViolation(err) {
			return errs.ErrAlreadyExists
		}
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := resolveQuerier(ctx, r.pool)

	var user domain.User
	if err := q.QueryRow(ctx, selectUserByEmailQuery, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUnauthorized
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	q := resolveQuerier(ctx, r.pool)

	var user domain.User
	if err := q.QueryRow(ctx, selectUserByIDQuery, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
