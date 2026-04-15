package auth

import (
	"context"
	"strings"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type loginUserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type loginPasswordHasher interface {
	Compare(hash, password string) error
}

type loginTokenService interface {
	GenerateAccessToken(userID int64) (string, error)
}

type LoginUseCase struct {
	users  loginUserRepository
	hasher loginPasswordHasher
	tokens loginTokenService
}

func NewLoginUseCase(users loginUserRepository, hasher loginPasswordHasher, tokens loginTokenService) *LoginUseCase {
	return &LoginUseCase{users: users, hasher: hasher, tokens: tokens}
}

func (uc *LoginUseCase) Execute(ctx context.Context, in dto.LoginInput) (*dto.AuthOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.Password == "" {
		return nil, errs.ErrInvalidInput
	}

	user, err := uc.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := uc.hasher.Compare(user.PasswordHash, in.Password); err != nil {
		return nil, errs.ErrUnauthorized
	}

	token, err := uc.tokens.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthOutput{AccessToken: token}, nil
}
