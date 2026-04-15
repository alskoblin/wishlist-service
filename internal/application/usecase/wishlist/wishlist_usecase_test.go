package wishlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type wishlistRepoMock struct {
	createFn          func(ctx context.Context, wishlist *domain.Wishlist) error
	listByOwnerFn     func(ctx context.Context, ownerID int64) ([]domain.Wishlist, error)
	getByIDAndOwnerFn func(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
	updateFn          func(ctx context.Context, wishlist *domain.Wishlist) error
	deleteFn          func(ctx context.Context, id, ownerID int64) error
	getByTokenFn      func(ctx context.Context, token string) (*domain.Wishlist, error)
}

func (m *wishlistRepoMock) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	if m.createFn != nil {
		return m.createFn(ctx, wishlist)
	}
	return nil
}

func (m *wishlistRepoMock) ListByOwner(ctx context.Context, ownerID int64) ([]domain.Wishlist, error) {
	if m.listByOwnerFn != nil {
		return m.listByOwnerFn(ctx, ownerID)
	}
	return nil, nil
}

func (m *wishlistRepoMock) GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error) {
	if m.getByIDAndOwnerFn != nil {
		return m.getByIDAndOwnerFn(ctx, id, ownerID)
	}
	return nil, nil
}

func (m *wishlistRepoMock) Update(ctx context.Context, wishlist *domain.Wishlist) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, wishlist)
	}
	return nil
}

func (m *wishlistRepoMock) Delete(ctx context.Context, id, ownerID int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id, ownerID)
	}
	return nil
}

func (m *wishlistRepoMock) GetByToken(ctx context.Context, token string) (*domain.Wishlist, error) {
	if m.getByTokenFn != nil {
		return m.getByTokenFn(ctx, token)
	}
	return nil, nil
}

type wishlistItemsMock struct {
	listByWishlistFn func(ctx context.Context, wishlistID int64) ([]domain.Item, error)
}

func (m *wishlistItemsMock) ListByWishlist(ctx context.Context, wishlistID int64) ([]domain.Item, error) {
	if m.listByWishlistFn != nil {
		return m.listByWishlistFn(ctx, wishlistID)
	}
	return nil, nil
}

type publicTokenMock struct {
	generatePublicTokenFn func() (string, error)
}

func (m *publicTokenMock) GeneratePublicToken() (string, error) {
	if m.generatePublicTokenFn != nil {
		return m.generatePublicTokenFn()
	}
	return "", nil
}

func TestCreateUseCaseExecuteSuccess(t *testing.T) {
	eventDate := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)
	var created domain.Wishlist

	repo := &wishlistRepoMock{
		createFn: func(_ context.Context, wishlist *domain.Wishlist) error {
			created = *wishlist
			return nil
		},
	}
	tokens := &publicTokenMock{generatePublicTokenFn: func() (string, error) { return "public-token", nil }}

	uc := NewCreateUseCase(repo, tokens)
	out, err := uc.Execute(context.Background(), dto.CreateWishlistInput{
		OwnerID:     1,
		EventTitle:  "  Birthday  ",
		Description: "  Gifts  ",
		EventDate:   eventDate,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.PublicToken != "public-token" {
		t.Fatalf("unexpected token: %s", out.PublicToken)
	}
	if created.EventTitle != "Birthday" {
		t.Fatalf("unexpected title: %s", created.EventTitle)
	}
	if created.Description != "Gifts" {
		t.Fatalf("unexpected description: %s", created.Description)
	}
}

func TestCreateUseCaseExecuteInvalidInput(t *testing.T) {
	uc := NewCreateUseCase(&wishlistRepoMock{}, &publicTokenMock{})

	_, err := uc.Execute(context.Background(), dto.CreateWishlistInput{OwnerID: 0, EventTitle: "Birthday"})
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}

func TestUpdateUseCaseExecuteSuccess(t *testing.T) {
	repo := &wishlistRepoMock{
		getByIDAndOwnerFn: func(_ context.Context, id, ownerID int64) (*domain.Wishlist, error) {
			return &domain.Wishlist{ID: id, OwnerID: ownerID, EventTitle: "Old"}, nil
		},
		updateFn: func(_ context.Context, wishlist *domain.Wishlist) error {
			if wishlist.EventTitle != "New title" {
				t.Fatalf("unexpected title in update: %s", wishlist.EventTitle)
			}
			return nil
		},
	}

	uc := NewUpdateUseCase(repo)
	out, err := uc.Execute(context.Background(), dto.UpdateWishlistInput{
		WishlistID: 10,
		OwnerID:    3,
		EventTitle: "  New title  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.EventTitle != "New title" {
		t.Fatalf("unexpected updated title: %s", out.EventTitle)
	}
}

func TestGetByTokenUseCaseExecuteSuccess(t *testing.T) {
	repo := &wishlistRepoMock{
		getByTokenFn: func(_ context.Context, token string) (*domain.Wishlist, error) {
			if token != "public-123" {
				t.Fatalf("unexpected token: %s", token)
			}
			return &domain.Wishlist{ID: 11, PublicToken: token}, nil
		},
	}
	itemsRepo := &wishlistItemsMock{
		listByWishlistFn: func(_ context.Context, wishlistID int64) ([]domain.Item, error) {
			if wishlistID != 11 {
				t.Fatalf("unexpected wishlistID: %d", wishlistID)
			}
			return []domain.Item{{ID: 1, Title: "Book"}}, nil
		},
	}

	uc := NewGetByTokenUseCase(repo, itemsRepo)
	out, err := uc.Execute(context.Background(), "public-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Wishlist.ID != 11 || len(out.Items) != 1 {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestGetByTokenUseCaseExecuteInvalidInput(t *testing.T) {
	uc := NewGetByTokenUseCase(&wishlistRepoMock{}, &wishlistItemsMock{})
	_, err := uc.Execute(context.Background(), "   ")
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
}
