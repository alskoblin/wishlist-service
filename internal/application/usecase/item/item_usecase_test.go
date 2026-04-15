package item

import (
	"context"
	"errors"
	"testing"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type itemWishlistRepoMock struct {
	getByIDAndOwnerFn func(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
}

func (m *itemWishlistRepoMock) GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error) {
	if m.getByIDAndOwnerFn != nil {
		return m.getByIDAndOwnerFn(ctx, id, ownerID)
	}
	return nil, nil
}

type itemRepoMock struct {
	createFn               func(ctx context.Context, item *domain.Item) error
	getByIDAndWishlistFn   func(ctx context.Context, itemID, wishlistID int64) (*domain.Item, error)
	updateFn               func(ctx context.Context, item *domain.Item) error
	deleteFn               func(ctx context.Context, itemID, wishlistID int64) error
	reserveByPublicTokenFn func(ctx context.Context, publicToken string, itemID int64) error
}

func (m *itemRepoMock) Create(ctx context.Context, item *domain.Item) error {
	if m.createFn != nil {
		return m.createFn(ctx, item)
	}
	return nil
}

func (m *itemRepoMock) GetByIDAndWishlist(ctx context.Context, itemID, wishlistID int64) (*domain.Item, error) {
	if m.getByIDAndWishlistFn != nil {
		return m.getByIDAndWishlistFn(ctx, itemID, wishlistID)
	}
	return nil, nil
}

func (m *itemRepoMock) Update(ctx context.Context, item *domain.Item) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, item)
	}
	return nil
}

func (m *itemRepoMock) Delete(ctx context.Context, itemID, wishlistID int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, itemID, wishlistID)
	}
	return nil
}

func (m *itemRepoMock) ReserveByPublicToken(ctx context.Context, publicToken string, itemID int64) error {
	if m.reserveByPublicTokenFn != nil {
		return m.reserveByPublicTokenFn(ctx, publicToken, itemID)
	}
	return nil
}

func TestReserveItemUseCaseExecuteSuccess(t *testing.T) {
	items := &itemRepoMock{
		reserveByPublicTokenFn: func(_ context.Context, token string, itemID int64) error {
			if token != "public-token" || itemID != 7 {
				t.Fatalf("unexpected args: token=%s itemID=%d", token, itemID)
			}
			return nil
		},
	}

	uc := NewReserveItemUseCase(items)
	if err := uc.Execute(context.Background(), "public-token", 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReserveItemUseCaseExecuteInvalidInput(t *testing.T) {
	uc := NewReserveItemUseCase(&itemRepoMock{})

	err := uc.Execute(context.Background(), "   ", 7)
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}

func TestCreateItemUseCaseExecuteSuccess(t *testing.T) {
	wishlists := &itemWishlistRepoMock{
		getByIDAndOwnerFn: func(_ context.Context, id, ownerID int64) (*domain.Wishlist, error) {
			return &domain.Wishlist{ID: id, OwnerID: ownerID}, nil
		},
	}
	items := &itemRepoMock{
		createFn: func(_ context.Context, item *domain.Item) error {
			if item.Title != "Book" {
				t.Fatalf("unexpected title: %s", item.Title)
			}
			if item.Description != "Interesting" {
				t.Fatalf("unexpected description: %s", item.Description)
			}
			return nil
		},
	}

	uc := NewCreateUseCase(wishlists, items)
	out, err := uc.Execute(context.Background(), dto.CreateItemInput{
		WishlistID:  10,
		OwnerID:     1,
		Title:       "  Book  ",
		Description: "  Interesting  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Title != "Book" {
		t.Fatalf("unexpected output title: %s", out.Title)
	}
}

func TestUpdateItemUseCaseExecuteSuccess(t *testing.T) {
	wishlists := &itemWishlistRepoMock{
		getByIDAndOwnerFn: func(_ context.Context, id, ownerID int64) (*domain.Wishlist, error) {
			return &domain.Wishlist{ID: id, OwnerID: ownerID}, nil
		},
	}
	items := &itemRepoMock{
		getByIDAndWishlistFn: func(_ context.Context, itemID, wishlistID int64) (*domain.Item, error) {
			return &domain.Item{ID: itemID, WishlistID: wishlistID, Title: "Old"}, nil
		},
		updateFn: func(_ context.Context, item *domain.Item) error {
			if item.Title != "New" {
				t.Fatalf("unexpected updated title: %s", item.Title)
			}
			return nil
		},
	}

	uc := NewUpdateUseCase(wishlists, items)
	out, err := uc.Execute(context.Background(), dto.UpdateItemInput{
		ItemID:     1,
		WishlistID: 2,
		OwnerID:    3,
		Title:      "  New  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Title != "New" {
		t.Fatalf("unexpected output title: %s", out.Title)
	}
}

func TestDeleteItemUseCaseExecuteInvalidInput(t *testing.T) {
	uc := NewDeleteUseCase(&itemWishlistRepoMock{}, &itemRepoMock{})
	err := uc.Execute(context.Background(), 0, 1, 1)
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}
