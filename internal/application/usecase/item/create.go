package item

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type createItemWishlistRepository interface {
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
}

type createItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
}

type CreateUseCase struct {
	wishlists createItemWishlistRepository
	items     createItemRepository
}

func NewCreateUseCase(wishlists createItemWishlistRepository, items createItemRepository) *CreateUseCase {
	return &CreateUseCase{wishlists: wishlists, items: items}
}

func (uc *CreateUseCase) Execute(ctx context.Context, in dto.CreateItemInput) (*domain.Item, error) {
	if in.WishlistID == 0 || in.OwnerID == 0 || strings.TrimSpace(in.Title) == "" {
		return nil, errs.ErrInvalidInput
	}

	if _, err := uc.wishlists.GetByIDAndOwner(ctx, in.WishlistID, in.OwnerID); err != nil {
		return nil, err
	}

	item := &domain.Item{
		WishlistID:  in.WishlistID,
		Title:       strings.TrimSpace(in.Title),
		Description: strings.TrimSpace(in.Description),
		ProductURL:  strings.TrimSpace(in.ProductURL),
		Priority:    in.Priority,
	}

	if err := uc.items.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}
