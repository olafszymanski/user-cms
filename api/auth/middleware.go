package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/olafszymanski/user-cms/postgres"
	"github.com/olafszymanski/user-cms/users"
)

type contextKey string

var ctxKey contextKey = "user"

type HttpError struct {
	Error string `json:"error"`
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		username, err := ParseToken(token)
		if err != nil {
			message := &HttpError{Error: err.Error()}
			jsonMessage, _ := json.Marshal(message)
			http.Error(w, string(jsonMessage), http.StatusForbidden)
			return
		}
		user, err := postgres.Database.UserStore.GetByUsername(*username)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(ctxKey).(*users.User)
	return raw
}
