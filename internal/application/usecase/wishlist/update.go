package wishlist

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type updateWishlistRepository interface {
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
	Update(ctx context.Context, wishlist *domain.Wishlist) error
}

type UpdateUseCase struct {
	wishlists updateWishlistRepository
}

func NewUpdateUseCase(wishlists updateWishlistRepository) *UpdateUseCase {
	return &UpdateUseCase{wishlists: wishlists}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, in dto.UpdateWishlistInput) (*domain.Wishlist, error) {
	if in.WishlistID == 0 || in.OwnerID == 0 || strings.TrimSpace(in.EventTitle) == "" {
		return nil, errs.ErrInvalidInput
	}

	wishlist, err := uc.wishlists.GetByIDAndOwner(ctx, in.WishlistID, in.OwnerID)
	if err != nil {
		return nil, err
	}

	wishlist.EventTitle = strings.TrimSpace(in.EventTitle)
	wishlist.Description = strings.TrimSpace(in.Description)
	wishlist.EventDate = in.EventDate

	if err := uc.wishlists.Update(ctx, wishlist); err != nil {
		return nil, err
	}

	return wishlist, nil
}
