package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/EricGusmao/easy-todo/user"
)

type contextKey string

const contextUserKey = contextKey("auth.user")

func UserFromContext(ctx context.Context) (*user.User, bool) {
	user, ok := ctx.Value(contextUserKey).(*user.User)
	if !ok || user == nil {
		return nil, false
	}
	return user, true
}

func NewContextWithUser(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, contextUserKey, user)
}

func NewMiddleware(auth Service) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authorization := r.Header.Get("Authorization")
				token := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))
				user, err := auth.UserFromToken(r.Context(), token)
				if err != nil || user == nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				rctx := NewContextWithUser(r.Context(), user)
				newReq := r.WithContext(rctx)
				h.ServeHTTP(w, newReq)
			},
		)
	}
}
