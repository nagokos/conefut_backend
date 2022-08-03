package auth

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/jwt"
	"github.com/nagokos/connefut_backend/graph/models/user"
)

func Middleware(dbPool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("token")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}
			if c.Value == "" {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				next.ServeHTTP(w, r)
				return
			}

			userID, err := jwt.ParseToken(c.Value)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				next.ServeHTTP(w, r)
				return
			}

			viewer, err := user.GetUser(r.Context(), dbPool, userID)
			if err != nil {
				http.Error(w, "user not found", http.StatusForbidden)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), user.UserCtxKey, viewer)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
