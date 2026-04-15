package handlers

import (
	"context"

	authuc "wishlist-service/internal/application/usecase/auth"
	itemuc "wishlist-service/internal/application/usecase/item"
	wishlistuc "wishlist-service/internal/application/usecase/wishlist"
	"wishlist-service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type WishlistRepository interface {
	Create(ctx context.Context, wishlist *domain.Wishlist) error
	ListByOwner(ctx context.Context, ownerID int64) ([]domain.Wishlist, error)
	GetByIDAndOwner(ctx context.Context, id, ownerID int64) (*domain.Wishlist, error)
	Update(ctx context.Context, wishlist *domain.Wishlist) error
	Delete(ctx context.Context, id, ownerID int64) error
	GetByToken(ctx context.Context, token string) (*domain.Wishlist, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item *domain.Item) error
	ListByWishlist(ctx context.Context, wishlistID int64) ([]domain.Item, error)
	GetByIDAndWishlist(ctx context.Context, itemID, wishlistID int64) (*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, itemID, wishlistID int64) error
	ReserveByPublicToken(ctx context.Context, publicToken string, itemID int64) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenService interface {
	GenerateAccessToken(userID int64) (string, error)
	GeneratePublicToken() (string, error)
}

type Dependencies struct {
	Users          UserRepository
	Wishlists      WishlistRepository
	Items          ItemRepository
	PasswordHasher PasswordHasher
	TokenService   TokenService
}

type Handler struct {
	registerUC *authuc.RegisterUseCase
	loginUC    *authuc.LoginUseCase

	createWishlistUC *wishlistuc.CreateUseCase
	listWishlistUC   *wishlistuc.ListUseCase
	updateWishlistUC *wishlistuc.UpdateUseCase
	deleteWishlistUC *wishlistuc.DeleteUseCase

	createItemUC *itemuc.CreateUseCase
	listItemUC   *itemuc.ListUseCase
	updateItemUC *itemuc.UpdateUseCase
	deleteItemUC *itemuc.DeleteUseCase

	publicGetUC     *wishlistuc.GetByTokenUseCase
	publicReserveUC *itemuc.ReserveItemUseCase
}

func NewHandler(dep Dependencies) *Handler {
	return &Handler{
		registerUC: authuc.NewRegisterUseCase(dep.Users, dep.PasswordHasher, dep.TokenService),
		loginUC:    authuc.NewLoginUseCase(dep.Users, dep.PasswordHasher, dep.TokenService),

		createWishlistUC: wishlistuc.NewCreateUseCase(dep.Wishlists, dep.TokenService),
		listWishlistUC:   wishlistuc.NewListUseCase(dep.Wishlists),
		updateWishlistUC: wishlistuc.NewUpdateUseCase(dep.Wishlists),
		deleteWishlistUC: wishlistuc.NewDeleteUseCase(dep.Wishlists),

		createItemUC: itemuc.NewCreateUseCase(dep.Wishlists, dep.Items),
		listItemUC:   itemuc.NewListUseCase(dep.Wishlists, dep.Items),
		updateItemUC: itemuc.NewUpdateUseCase(dep.Wishlists, dep.Items),
		deleteItemUC: itemuc.NewDeleteUseCase(dep.Wishlists, dep.Items),

		publicGetUC:     wishlistuc.NewGetByTokenUseCase(dep.Wishlists, dep.Items),
		publicReserveUC: itemuc.NewReserveItemUseCase(dep.Items),
	}
}
