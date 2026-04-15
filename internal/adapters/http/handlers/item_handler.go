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

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
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

	items, err := h.listItemUC.Execute(r.Context(), wishlistID, userID)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
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

	var in dto.CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}
	in.OwnerID = userID
	in.WishlistID = wishlistID

	item, err := h.createItemUC.Execute(r.Context(), in)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
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

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	var in dto.UpdateItemInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}
	in.OwnerID = userID
	in.WishlistID = wishlistID
	in.ItemID = itemID

	item, err := h.updateItemUC.Execute(r.Context(), in)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
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

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.deleteItemUC.Execute(r.Context(), itemID, wishlistID, userID); err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
