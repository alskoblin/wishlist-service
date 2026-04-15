package wishlist

import (
	"context"

	"wishlist-service/internal/errs"
)

type deleteWishlistRepository interface {
	Delete(ctx context.Context, id, ownerID int64) error
}

type DeleteUseCase struct {
	wishlists deleteWishlistRepository
}

func NewDeleteUseCase(wishlists deleteWishlistRepository) *DeleteUseCase {
	return &DeleteUseCase{wishlists: wishlists}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, wishlistID, ownerID int64) error {
	if wishlistID == 0 || ownerID == 0 {
		return errs.ErrInvalidInput
	}

	return uc.wishlists.Delete(ctx, wishlistID, ownerID)
}
