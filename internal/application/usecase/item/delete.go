package item

import (
	"context"

	"wishlist-service/internal/errs"
)

type DeleteUseCase struct {
	wishlists WishlistRepository
	items     ItemRepository
}

func NewDeleteUseCase(wishlists WishlistRepository, items ItemRepository) *DeleteUseCase {
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
