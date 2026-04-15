package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"wishlist-service/internal/adapters/http/middleware"
	"wishlist-service/internal/adapters/http/presenter"
	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/errs"
)

func (h *Handler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		presenter.WriteError(w, errs.ErrUnauthorized)
		return
	}

	var in dto.CreateWishlistInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}
	in.OwnerID = userID

	wishlist, err := h.createWishlistUC.Execute(r.Context(), in)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusCreated, wishlist)
}

func (h *Handler) ListWishlists(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		presenter.WriteError(w, errs.ErrUnauthorized)
		return
	}

	wishlists, err := h.listWishlistUC.Execute(r.Context(), userID)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, wishlists)
}

func (h *Handler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		presenter.WriteError(w, errs.ErrUnauthorized)
		return
	}

	wishlistID, err := strconv.ParseInt(chi.URLParam(r, "wishlistID"), 10, 64)
	if err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	var in dto.UpdateWishlistInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}
	in.WishlistID = wishlistID
	in.OwnerID = userID

	wishlist, err := h.updateWishlistUC.Execute(r.Context(), in)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, wishlist)
}

func (h *Handler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		presenter.WriteError(w, errs.ErrUnauthorized)
		return
	}

	wishlistID, err := strconv.ParseInt(chi.URLParam(r, "wishlistID"), 10, 64)
	if err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.deleteWishlistUC.Execute(r.Context(), wishlistID, userID); err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
