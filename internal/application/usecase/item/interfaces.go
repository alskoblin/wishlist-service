package item

import (
	"context"

	"wishlist-service/internal/domain"
)

type WishlistRepository interface {
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	GetByIDAndWishlist(ctx context.Context, itemID, wishlistID int64) (*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, itemID, wishlistID int64) error
	ReserveByPublicToken(ctx context.Context, publicToken string, itemID int64) error
}
