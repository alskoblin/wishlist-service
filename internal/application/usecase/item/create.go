package item

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type CreateUseCase struct {
	wishlists WishlistRepository
	items     ItemRepository
}

func NewCreateUseCase(wishlists WishlistRepository, items ItemRepository) *CreateUseCase {
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
