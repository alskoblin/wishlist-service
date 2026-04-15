package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"wishlist-service/internal/adapters/http/presenter"
	"wishlist-service/internal/errs"
)

func (h *Handler) GetPublicWishlist(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	out, err := h.publicGetUC.Execute(r.Context(), token)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) ReservePublicItem(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	if err := h.publicReserveUC.Execute(r.Context(), token, itemID); err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, map[string]string{"status": "reserved"})
}
