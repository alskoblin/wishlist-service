package item

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type updateItemWishlistRepository interface {
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
}

type updateItemRepository interface {
	GetByIDAndWishlist(ctx context.Context, itemID, wishlistID int64) (*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
}

type UpdateUseCase struct {
	wishlists updateItemWishlistRepository
	items     updateItemRepository
}

func NewUpdateUseCase(wishlists updateItemWishlistRepository, items updateItemRepository) *UpdateUseCase {
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
