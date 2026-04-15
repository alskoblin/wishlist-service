package wishlist

import (
	"context"
	"strings"

	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type getByTokenWishlistRepository interface {
	GetByToken(ctx context.Context, token string) (*domain.Wishlist, error)
}

type getByTokenItemRepository interface {
	ListByWishlist(ctx context.Context, wishlistID int64) ([]domain.Item, error)
}

type GetByTokenOutput struct {
	Wishlist *domain.Wishlist `json:"wishlist"`
	Items    []domain.Item    `json:"items"`
}

type GetByTokenUseCase struct {
	wishlists getByTokenWishlistRepository
	items     getByTokenItemRepository
}

func NewGetByTokenUseCase(wishlists getByTokenWishlistRepository, items getByTokenItemRepository) *GetByTokenUseCase {
	return &GetByTokenUseCase{wishlists: wishlists, items: items}
}

func (uc *GetByTokenUseCase) Execute(ctx context.Context, token string) (*GetByTokenOutput, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errs.ErrInvalidInput
	}

	wishlist, err := uc.wishlists.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	items, err := uc.items.ListByWishlist(ctx, wishlist.ID)
	if err != nil {
		return nil, err
	}

	return &GetByTokenOutput{Wishlist: wishlist, Items: items}, nil
}
