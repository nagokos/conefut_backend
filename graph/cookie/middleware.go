package cookie

import (
	"context"
	"net/http"
	"time"
)

var httpWriterKey = &contextKey{name: "httpWriter"}

type contextKey struct {
	name string
}

func MiddleWare() func(http.Handler) http.Handler {
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
		Name:     "token",
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
		Name:     "token",
	})
}
