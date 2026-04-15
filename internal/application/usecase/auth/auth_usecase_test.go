package auth

import (
	"context"
	"errors"
	"testing"

	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/domain"
	"wishlist-service/internal/errs"
)

type authUserRepoMock struct {
	createFn     func(ctx context.Context, user *domain.User) error
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
}

func (m *authUserRepoMock) Create(ctx context.Context, user *domain.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return nil
}

func (m *authUserRepoMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}

type authHasherMock struct {
	hashFn    func(password string) (string, error)
	compareFn func(hash, password string) error
}

func (m *authHasherMock) Hash(password string) (string, error) {
	if m.hashFn != nil {
		return m.hashFn(password)
	}
	return "", nil
}

func (m *authHasherMock) Compare(hash, password string) error {
	if m.compareFn != nil {
		return m.compareFn(hash, password)
	}
	return nil
}

type authTokenMock struct {
	generateAccessTokenFn func(userID int64) (string, error)
}

func (m *authTokenMock) GenerateAccessToken(userID int64) (string, error) {
	if m.generateAccessTokenFn != nil {
		return m.generateAccessTokenFn(userID)
	}
	return "", nil
}

func TestRegisterUseCaseExecuteSuccess(t *testing.T) {
	ctx := context.Background()
	var createdUser domain.User

	repo := &authUserRepoMock{
		createFn: func(_ context.Context, user *domain.User) error {
			createdUser = *user
			user.ID = 42
			return nil
		},
	}
	hasher := &authHasherMock{
		hashFn: func(password string) (string, error) {
			if password != "password123" {
				t.Fatalf("unexpected password: %s", password)
			}
			return "hashed-pass", nil
		},
	}
	tokens := &authTokenMock{
		generateAccessTokenFn: func(userID int64) (string, error) {
			if userID != 42 {
				t.Fatalf("unexpected userID: %d", userID)
			}
			return "jwt-token", nil
		},
	}

	uc := NewRegisterUseCase(repo, hasher, tokens)
	out, err := uc.Execute(ctx, dto.RegisterInput{Email: "  USER@Example.COM  ", Password: "password123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.AccessToken != "jwt-token" {
		t.Fatalf("unexpected access token: %s", out.AccessToken)
	}
	if createdUser.Email != "user@example.com" {
		t.Fatalf("unexpected normalized email: %s", createdUser.Email)
	}
	if createdUser.PasswordHash != "hashed-pass" {
		t.Fatalf("unexpected password hash: %s", createdUser.PasswordHash)
	}
}

func TestRegisterUseCaseExecuteInvalidInput(t *testing.T) {
	uc := NewRegisterUseCase(&authUserRepoMock{}, &authHasherMock{}, &authTokenMock{})

	_, err := uc.Execute(context.Background(), dto.RegisterInput{Email: "", Password: "12345678"})
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}

	_, err = uc.Execute(context.Background(), dto.RegisterInput{Email: "user@example.com", Password: "short"})
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for short password, got: %v", err)
	}
}

func TestLoginUseCaseExecuteSuccess(t *testing.T) {
	repo := &authUserRepoMock{
		getByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			if email != "user@example.com" {
				t.Fatalf("unexpected email: %s", email)
			}
			return &domain.User{ID: 99, Email: email, PasswordHash: "stored-hash"}, nil
		},
	}
	hasher := &authHasherMock{
		compareFn: func(hash, password string) error {
			if hash != "stored-hash" || password != "password123" {
				t.Fatalf("unexpected compare args: hash=%s password=%s", hash, password)
			}
			return nil
		},
	}
	tokens := &authTokenMock{
		generateAccessTokenFn: func(userID int64) (string, error) {
			if userID != 99 {
				t.Fatalf("unexpected userID: %d", userID)
			}
			return "login-token", nil
		},
	}

	uc := NewLoginUseCase(repo, hasher, tokens)
	out, err := uc.Execute(context.Background(), dto.LoginInput{Email: " USER@EXAMPLE.COM ", Password: "password123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.AccessToken != "login-token" {
		t.Fatalf("unexpected access token: %s", out.AccessToken)
	}
}

func TestLoginUseCaseExecuteWrongPassword(t *testing.T) {
	repo := &authUserRepoMock{
		getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: "user@example.com", PasswordHash: "stored-hash"}, nil
		},
	}
	hasher := &authHasherMock{compareFn: func(_, _ string) error { return errors.New("mismatch") }}
	uc := NewLoginUseCase(repo, hasher, &authTokenMock{})

	_, err := uc.Execute(context.Background(), dto.LoginInput{Email: "user@example.com", Password: "wrong"})
	if !errors.Is(err, errs.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got: %v", err)
	}
}
