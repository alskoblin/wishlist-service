package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"wishlist-service/docs/openapi"
	"wishlist-service/internal/adapters/http/handlers"
	appmw "wishlist-service/internal/adapters/http/middleware"
)

func NewRouter(dep handlers.Dependencies, tokenService appmw.AccessTokenParser) http.Handler {
	h := handlers.NewHandler(dep)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(appmw.Recoverer)

	r.Get("/openapi.yaml", openapi.SpecHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)

		r.Get("/public/{token}", h.GetPublicWishlist)
		r.Post("/public/{token}/reserve/{itemID}", h.ReservePublicItem)

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthRequired(tokenService))

			r.Route("/wishlists", func(r chi.Router) {
				r.Post("/", h.CreateWishlist)
				r.Get("/", h.ListWishlists)
				r.Put("/{wishlistID}", h.UpdateWishlist)
				r.Delete("/{wishlistID}", h.DeleteWishlist)

				r.Route("/{wishlistID}/items", func(r chi.Router) {
					r.Post("/", h.CreateItem)
					r.Put("/{itemID}", h.UpdateItem)
					r.Delete("/{itemID}", h.DeleteItem)
				})
			})
		})
	})

	return r
}
