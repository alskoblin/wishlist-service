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

type ItemRepository struct {
	pool *pgxpool.Pool
}

func NewItemRepository(pool *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{pool: pool}
}

const (
	insertItemQuery = `
		INSERT INTO wishlist_items (wishlist_id, title, description, product_url, priority)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, reserved, reserved_at, created_at, updated_at
	`
	selectItemsByWishlistQuery = `
		SELECT id, wishlist_id, title, description, product_url, priority, reserved, reserved_at, created_at, updated_at
		FROM wishlist_items
		WHERE wishlist_id = $1
		ORDER BY priority DESC, created_at ASC
	`
	selectItemByIDAndWishlistQuery = `
		SELECT id, wishlist_id, title, description, product_url, priority, reserved, reserved_at, created_at, updated_at
		FROM wishlist_items
		WHERE id = $1 AND wishlist_id = $2
	`
	updateItemQuery = `
		UPDATE wishlist_items
		SET title = $1,
			description = $2,
			product_url = $3,
			priority = $4,
			reserved = $5,
			reserved_at = $6,
			updated_at = NOW()
		WHERE id = $7 AND wishlist_id = $8
		RETURNING updated_at
	`
	deleteItemQuery = `
		DELETE FROM wishlist_items
		WHERE id = $1 AND wishlist_id = $2
	`
	reserveItemByPublicTokenQuery = `
		UPDATE wishlist_items wi
		SET reserved = TRUE,
			reserved_at = NOW(),
			updated_at = NOW()
		FROM wishlists w
		WHERE wi.id = $1
		  AND wi.wishlist_id = w.id
		  AND w.public_token = $2
		  AND wi.reserved = FALSE
		RETURNING wi.id
	`
	selectReservedStateByPublicTokenQuery = `
		SELECT wi.reserved
		FROM wishlist_items wi
		JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE wi.id = $1
		  AND w.public_token = $2
	`
)

func (r *ItemRepository) Create(ctx context.Context, item *domain.Item) error {
	q := resolveQuerier(ctx, r.pool)

	if err := q.QueryRow(ctx, insertItemQuery, item.WishlistID, item.Title, item.Description, item.ProductURL, item.Priority).
		Scan(&item.ID, &item.Reserved, &item.ReservedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
		return fmt.Errorf("insert item: %w", err)
	}

	return nil
}

func (r *ItemRepository) ListByWishlist(ctx context.Context, wishlistID int64) ([]domain.Item, error) {
	q := resolveQuerier(ctx, r.pool)

	rows, err := q.Query(ctx, selectItemsByWishlistQuery, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Item, 0)
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.WishlistID, &item.Title, &item.Description, &item.ProductURL, &item.Priority, &item.Reserved, &item.ReservedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate items: %w", err)
	}

	return items, nil
}

func (r *ItemRepository) GetByIDAndWishlist(ctx context.Context, itemID, wishlistID int64) (*domain.Item, error) {
	q := resolveQuerier(ctx, r.pool)

	var item domain.Item
	if err := q.QueryRow(ctx, selectItemByIDAndWishlistQuery, itemID, wishlistID).Scan(&item.ID, &item.WishlistID, &item.Title, &item.Description, &item.ProductURL, &item.Priority, &item.Reserved, &item.ReservedAt, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("get item: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) Update(ctx context.Context, item *domain.Item) error {
	q := resolveQuerier(ctx, r.pool)

	if err := q.QueryRow(ctx, updateItemQuery, item.Title, item.Description, item.ProductURL, item.Priority, item.Reserved, item.ReservedAt, item.ID, item.WishlistID).
		Scan(&item.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrNotFound
		}
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, itemID, wishlistID int64) error {
	q := resolveQuerier(ctx, r.pool)

	res, err := q.Exec(ctx, deleteItemQuery, itemID, wishlistID)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	if res.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *ItemRepository) ReserveByPublicToken(ctx context.Context, publicToken string, itemID int64) error {
	q := resolveQuerier(ctx, r.pool)

	var reservedItemID int64
	if err := q.QueryRow(ctx, reserveItemByPublicTokenQuery, itemID, publicToken).Scan(&reservedItemID); err == nil {
		return nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("reserve item by token: %w", err)
	}

	var reserved bool
	if err := q.QueryRow(ctx, selectReservedStateByPublicTokenQuery, itemID, publicToken).Scan(&reserved); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrNotFound
		}
		return fmt.Errorf("check reserve state by token: %w", err)
	}

	if reserved {
		return errs.ErrAlreadyReserved
	}

	return errs.ErrNotFound
}
