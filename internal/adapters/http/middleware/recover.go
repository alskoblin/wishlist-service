package middleware

import (
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func Recoverer(next http.Handler) http.Handler {
	return chimw.Recoverer(next)
}
