package wishlist

import (
	"context"

	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type listWishlistRepository interface {
	ListByOwner(ctx context.Context, ownerID int64) ([]domain.Wishlist, error)
}

type ListUseCase struct {
	wishlists listWishlistRepository
}

func NewListUseCase(wishlists listWishlistRepository) *ListUseCase {
	return &ListUseCase{wishlists: wishlists}
}

func (uc *ListUseCase) Execute(ctx context.Context, ownerID int64) ([]domain.Wishlist, error) {
	if ownerID == 0 {
		return nil, errs.ErrInvalidInput
	}

	return uc.wishlists.ListByOwner(ctx, ownerID)
}
