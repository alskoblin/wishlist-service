package item

import (
	"context"

	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type listItemWishlistRepository interface {
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
}

type listItemRepository interface {
	ListByWishlist(ctx context.Context, wishlistID int64) ([]domain.Item, error)
}

type ListUseCase struct {
	wishlists listItemWishlistRepository
	items     listItemRepository
}

func NewListUseCase(wishlists listItemWishlistRepository, items listItemRepository) *ListUseCase {
	return &ListUseCase{wishlists: wishlists, items: items}
}

func (uc *ListUseCase) Execute(ctx context.Context, wishlistID, ownerID int64) ([]domain.Item, error) {
	if wishlistID == 0 || ownerID == 0 {
		return nil, errs.ErrInvalidInput
	}

	if _, err := uc.wishlists.GetByIDAndOwner(ctx, wishlistID, ownerID); err != nil {
		return nil, err
	}

	return uc.items.ListByWishlist(ctx, wishlistID)
}
