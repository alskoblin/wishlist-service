package auth

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type registerUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
}

type registerPasswordHasher interface {
	Hash(password string) (string, error)
}

type registerTokenService interface {
	GenerateAccessToken(userID int64) (string, error)
}

type RegisterUseCase struct {
	users  registerUserRepository
	hasher registerPasswordHasher
	tokens registerTokenService
}

func NewRegisterUseCase(users registerUserRepository, hasher registerPasswordHasher, tokens registerTokenService) *RegisterUseCase {
	return &RegisterUseCase{users: users, hasher: hasher, tokens: tokens}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, in dto.RegisterInput) (*dto.AuthOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || len(in.Password) < 8 {
		return nil, errs.ErrInvalidInput
	}

	hash, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{Email: email, PasswordHash: hash}
	if err := uc.users.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := uc.tokens.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthOutput{AccessToken: token}, nil
}
