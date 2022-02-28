package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/graph/domain/user"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

var userCtxKey = &contextKey{name: "secret"}
var httpWriterKey = &contextKey{name: "httpWriter"}

type contextKey struct {
	name string
}

func Middleware(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("jwt")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			userId, err := validateAndGetUserID(c)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				next.ServeHTTP(w, r)
				return
			}

			user, err := getUserByID(db, userId)
			if err != nil {
				http.Error(w, "user not found", http.StatusForbidden)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

func validateAndGetUserID(c *http.Cookie) (userID string, err error) {
	if time.Now().Before(c.Expires) {
		logger.Log.Error().Msg("expired...")
		return
	}

	if c.Value == "" {
		logger.Log.Error().Msg("jwt emypty!")
		return
	}

	token, err := jwt.Parse(c.Value, func(t *jwt.Token) (interface{}, error) {
		return user.SecretKey, nil
	})

	if err != nil {
		logger.Log.Error().Msg(err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	userID = claims["user_id"].(string)
	if userID == "" {
		logger.Log.Error().Msg("user_id not found")
		return
	}

	return userID, nil
}

func getUserByID(db *ent.Client, id string) (*model.User, error) {
	ctx := context.Background()
	res, err := db.User.Get(ctx, id)
	user := &model.User{
		ID:                      res.ID,
		Name:                    res.Name,
		Email:                   res.Email,
		Role:                    model.Role(res.Role),
		Avatar:                  res.Avatar,
		Introduction:            &res.Introduction,
		EmailVerificationStatus: model.EmailVerificationStatus(res.EmailVerificationStatus),
	}
	return user, err
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
		Name:     "jwt",
	})
}
