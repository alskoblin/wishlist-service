package item

import (
	"context"
	"strings"

	"wishlist-service/internal/errs"
)

type ReserveItemUseCase struct {
	items ItemRepository
}

func NewReserveItemUseCase(items ItemRepository) *ReserveItemUseCase {
	return &ReserveItemUseCase{items: items}
}

func (uc *ReserveItemUseCase) Execute(ctx context.Context, token string, itemID int64) error {
	if strings.TrimSpace(token) == "" || itemID == 0 {
		return errs.ErrInvalidInput
	}

	return uc.items.ReserveByPublicToken(ctx, token, itemID)
}
