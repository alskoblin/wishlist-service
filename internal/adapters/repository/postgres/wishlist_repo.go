package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type WishlistRepository struct {
	pool *pgxpool.Pool
}

func NewWishlistRepository(pool *pgxpool.Pool) *WishlistRepository {
	return &WishlistRepository{pool: pool}
}

const (
	insertWishlistQuery = `
		INSERT INTO wishlists (owner_id, event_title, description, event_date, public_token)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	selectWishlistsByOwnerQuery = `
		SELECT id, owner_id, event_title, description, event_date, public_token, created_at, updated_at
		FROM wishlists
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	selectWishlistByIDAndOwnerQuery = `
		SELECT id, owner_id, event_title, description, event_date, public_token, created_at, updated_at
		FROM wishlists
		WHERE id = $1 AND owner_id = $2
	`
	updateWishlistQuery = `
		UPDATE wishlists
		SET event_title = $1,
			description = $2,
			event_date = $3,
			updated_at = NOW()
		WHERE id = $4 AND owner_id = $5
		RETURNING updated_at
	`
	deleteWishlistQuery = `
		DELETE FROM wishlists
		WHERE id = $1 AND owner_id = $2
	`
	selectWishlistByTokenQuery = `
		SELECT id, owner_id, event_title, description, event_date, public_token, created_at, updated_at
		FROM wishlists
		WHERE public_token = $1
	`
)

func (r *WishlistRepository) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	q := resolveQuerier(ctx, r.pool)

	if err := q.QueryRow(ctx, insertWishlistQuery, wishlist.OwnerID, wishlist.EventTitle, wishlist.Description, wishlist.EventDate, wishlist.PublicToken).
		Scan(&wishlist.ID, &wishlist.CreatedAt, &wishlist.UpdatedAt); err != nil {
		if isUniqueViolation(err) {
			return errs.ErrAlreadyExists
		}
		return fmt.Errorf("insert wishlist: %w", err)
	}

	return nil
}

func (r *WishlistRepository) ListByOwner(ctx context.Context, ownerID int64) ([]domain.Wishlist, error) {
	q := resolveQuerier(ctx, r.pool)

	rows, err := q.Query(ctx, selectWishlistsByOwnerQuery, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list wishlists: %w", err)
	}
	defer rows.Close()

	wishlists := make([]domain.Wishlist, 0)
	for rows.Next() {
		var w domain.Wishlist
		if err := rows.Scan(&w.ID, &w.OwnerID, &w.EventTitle, &w.Description, &w.EventDate, &w.PublicToken, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan wishlist: %w", err)
		}
		wishlists = append(wishlists, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wishlists: %w", err)
	}

	return wishlists, nil
}

func (r *WishlistRepository) GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error) {
	q := resolveQuerier(ctx, r.pool)

	var w domain.Wishlist
	if err := q.QueryRow(ctx, selectWishlistByIDAndOwnerQuery, id, ownerID).Scan(&w.ID, &w.OwnerID, &w.EventTitle, &w.Description, &w.EventDate, &w.PublicToken, &w.CreatedAt, &w.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("get wishlist by id and owner: %w", err)
	}

	return &w, nil
}

func (r *WishlistRepository) Update(ctx context.Context, wishlist *domain.Wishlist) error {
	q := resolveQuerier(ctx, r.pool)

	if err := q.QueryRow(ctx, updateWishlistQuery, wishlist.EventTitle, wishlist.Description, wishlist.EventDate, wishlist.ID, wishlist.OwnerID).
		Scan(&wishlist.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrNotFound
		}
		return fmt.Errorf("update wishlist: %w", err)
	}

	return nil
}

func (r *WishlistRepository) Delete(ctx context.Context, id, ownerID int64) error {
	q := resolveQuerier(ctx, r.pool)

	res, err := q.Exec(ctx, deleteWishlistQuery, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete wishlist: %w", err)
	}
	if res.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *WishlistRepository) GetByToken(ctx context.Context, token string) (*domain.Wishlist, error) {
	q := resolveQuerier(ctx, r.pool)

	var w domain.Wishlist
	if err := q.QueryRow(ctx, selectWishlistByTokenQuery, token).Scan(&w.ID, &w.OwnerID, &w.EventTitle, &w.Description, &w.EventDate, &w.PublicToken, &w.CreatedAt, &w.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("get wishlist by token: %w", err)
	}

	return &w, nil
}
