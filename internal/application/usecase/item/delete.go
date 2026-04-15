package item

import (
	"context"

	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type deleteItemWishlistRepository interface {
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
}

type deleteItemRepository interface {
	Delete(ctx context.Context, itemID, wishlistID int64) error
}

type DeleteUseCase struct {
	wishlists deleteItemWishlistRepository
	items     deleteItemRepository
}

func NewDeleteUseCase(wishlists deleteItemWishlistRepository, items deleteItemRepository) *DeleteUseCase {
	return &DeleteUseCase{wishlists: wishlists, items: items}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, itemID, wishlistID, ownerID int64) error {
	if itemID == 0 || wishlistID == 0 || ownerID == 0 {
		return errs.ErrInvalidInput
	}

	if _, err := uc.wishlists.GetByIDAndOwner(ctx, wishlistID, ownerID); err != nil {
		return err
	}

	return uc.items.Delete(ctx, itemID, wishlistID)
}
