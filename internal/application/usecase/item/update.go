package item

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type UpdateUseCase struct {
	wishlists WishlistRepository
	items     ItemRepository
}

func NewUpdateUseCase(wishlists WishlistRepository, items ItemRepository) *UpdateUseCase {
	return &UpdateUseCase{wishlists: wishlists, items: items}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, in dto.UpdateItemInput) (*domain.Item, error) {
	if in.ItemID == 0 || in.WishlistID == 0 || in.OwnerID == 0 || strings.TrimSpace(in.Title) == "" {
		return nil, errs.ErrInvalidInput
	}

	if _, err := uc.wishlists.GetByIDAndOwner(ctx, in.WishlistID, in.OwnerID); err != nil {
		return nil, err
	}

	item, err := uc.items.GetByIDAndWishlist(ctx, in.ItemID, in.WishlistID)
	if err != nil {
		return nil, err
	}

	item.Title = strings.TrimSpace(in.Title)
	item.Description = strings.TrimSpace(in.Description)
	item.ProductURL = strings.TrimSpace(in.ProductURL)
	item.Priority = in.Priority

	if err := uc.items.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}
