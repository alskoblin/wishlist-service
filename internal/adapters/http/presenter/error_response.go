package presenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"wishlist-service/internal/errs"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, errs.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, errs.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, errs.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, errs.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, errs.ErrAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, errs.ErrAlreadyReserved):
		status = http.StatusConflict
	}

	writeJSON(w, status, ErrorResponse{Error: err.Error()})
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	writeJSON(w, status, payload)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
