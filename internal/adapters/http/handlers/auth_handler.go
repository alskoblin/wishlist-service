package handlers

import (
	"encoding/json"
	"net/http"

	"wishlist-service/internal/adapters/http/presenter"
	"wishlist-service/internal/application/dto"
	"wishlist-service/internal/errs"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in dto.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	out, err := h.registerUC.Execute(r.Context(), in)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusCreated, out)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in dto.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.WriteError(w, errs.ErrInvalidInput)
		return
	}

	out, err := h.loginUC.Execute(r.Context(), in)
	if err != nil {
		presenter.WriteError(w, err)
		return
	}

	presenter.WriteJSON(w, http.StatusOK, out)
}
