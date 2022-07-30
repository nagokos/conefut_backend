package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

var httpWriterKey = &contextKey{name: "httpWriter"}

type contextKey struct {
	name string
}

func Middleware(dbPool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("jwt")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			userID, err := validateAndGetUserID(c)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				next.ServeHTTP(w, r)
				return
			}

			viewer, err := user.GetUser(r.Context(), dbPool, utils.GenerateUniqueID("User", int(userID)))
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

func validateAndGetUserID(c *http.Cookie) (float64, error) {
	if time.Now().Before(c.Expires) {
		logger.NewLogger().Error("expired...")
		return 0, errors.New("expired")
	}

	if c.Value == "" {
		logger.NewLogger().Error("jwt empty!")
		return 0, errors.New("jwt empty")
	}

	token, err := jwt.Parse(c.Value, func(t *jwt.Token) (interface{}, error) {
		return user.SecretKey, nil
	})

	if err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	viewerID := claims["user_id"].(float64)

	return viewerID, nil
}

func CookieMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), httpWriterKey, w)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func SetAuthCookie(ctx context.Context, token string) {
	w, _ := ctx.Value(httpWriterKey).(http.ResponseWriter)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})
}

func RemoveAuthCookie(ctx context.Context) {
	w, _ := ctx.Value(httpWriterKey).(http.ResponseWriter)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, -1),
		Path:     "/",
		Name:     "jwt",
	})
}
