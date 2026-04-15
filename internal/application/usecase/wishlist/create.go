package wishlist

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type createWishlistRepository interface {
	Create(ctx context.Context, wishlist *domain.Wishlist) error
}

type createTokenService interface {
	GeneratePublicToken() (string, error)
}

type CreateUseCase struct {
	wishlists createWishlistRepository
	tokens    createTokenService
}

func NewCreateUseCase(wishlists createWishlistRepository, tokens createTokenService) *CreateUseCase {
	return &CreateUseCase{wishlists: wishlists, tokens: tokens}
}

func (uc *CreateUseCase) Execute(ctx context.Context, in dto.CreateWishlistInput) (*domain.Wishlist, error) {
	if in.OwnerID == 0 || strings.TrimSpace(in.EventTitle) == "" {
		return nil, errs.ErrInvalidInput
	}

	token, err := uc.tokens.GeneratePublicToken()
	if err != nil {
		return nil, err
	}

	wishlist := &domain.Wishlist{
		OwnerID:     in.OwnerID,
		EventTitle:  strings.TrimSpace(in.EventTitle),
		Description: strings.TrimSpace(in.Description),
		EventDate:   in.EventDate,
		PublicToken: token,
	}

	if err := uc.wishlists.Create(ctx, wishlist); err != nil {
		return nil, err
	}

	return wishlist, nil
}
