package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	httpadapter "wishlist-service/internal/adapters/http"
	"wishlist-service/internal/adapters/http/handlers"
	"wishlist-service/internal/adapters/repository/postgres"
	"wishlist-service/internal/infrastructure/auth"
	"wishlist-service/internal/infrastructure/config"
	"wishlist-service/internal/infrastructure/db"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	pool   *pgxpool.Pool
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if err := db.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	passwordHasher := auth.NewBcryptHasher()
	tokenService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTTTL)

	deps := handlers.Dependencies{
		Users:          postgres.NewUserRepository(pool),
		Wishlists:      postgres.NewWishlistRepository(pool),
		Items:          postgres.NewItemRepository(pool),
		PasswordHasher: passwordHasher,
		TokenService:   tokenService,
	}

	router := httpadapter.NewRouter(deps, tokenService)
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		cfg:    cfg,
		server: server,
		pool:   pool,
	}, nil
}

func (a *App) Server() *http.Server {
	return a.server
}

func (a *App) HTTPAddr() string {
	return a.cfg.HTTPAddr
}

func (a *App) Close() {
	if a.pool != nil {
		a.pool.Close()
	}
}
