package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/logger"
)

var userCtxKey = &contextKey{name: "secret"}
var httpWriterKey = &contextKey{name: "httpWriter"}

type contextKey struct {
	name string
}

func Middleware(dbConnection *sql.DB) func(http.Handler) http.Handler {
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

			user, err := getUserByID(dbConnection, userId)
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

func validateAndGetUserID(c *http.Cookie) (string, error) {
	if time.Now().Before(c.Expires) {
		logger.NewLogger().Error("expired...")
		return "", nil
	}

	if c.Value == "" {
		logger.NewLogger().Error("jwt empty!")
		return "", nil
	}

	token, err := jwt.Parse(c.Value, func(t *jwt.Token) (interface{}, error) {
		return user.SecretKey, nil
	})

	if err != nil {
		logger.NewLogger().Error(err.Error())
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.NewLogger().Error("Couldt't parse claims")
		return "", errors.New("Couldt't parse claims")
	}

	userID := claims["user_id"].(string)
	if userID == "" {
		logger.NewLogger().Error("user_id not found")
		return "", errors.New("user_id not found")
	}

	return userID, nil
}

func getUserByID(dbConnection *sql.DB, ID string) (*model.User, error) {
	var u model.User

	cmd := fmt.Sprintf("SELECT id, name, email, role, avatar, introduction, email_verification_status FROM %s WHERE id = $1", db.UserTable)
	row := dbConnection.QueryRow(cmd, ID)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Avatar, &u.Introduction, &u.EmailVerificationStatus)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &u, nil
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
