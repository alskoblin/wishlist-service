package middleware

import (
	"context"
	"net/http"
	"strings"

	"wishlist-service/internal/adapters/http/presenter"
	"wishlist-service/internal/errs"
)

type userIDCtxKey struct{}

type AccessTokenParser interface {
	ParseAccessToken(token string) (int64, error)
}

func AuthRequired(tokens AccessTokenParser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := parseBearerToken(r.Header.Get("Authorization"))
			if !ok {
				writeUnauthorized(w)
				return
			}

			userID, err := tokens.ParseAccessToken(token)
			if err != nil {
				writeUnauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), userIDCtxKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDCtxKey{}).(int64)
	return userID, ok
}

func parseBearerToken(header string) (string, bool) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", false
	}

	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}

	return token, true
}

func writeUnauthorized(w http.ResponseWriter) {
	presenter.WriteError(w, errs.ErrUnauthorized)
}
