package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/models/authentication"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/oauth2"
)

func lineProvider(ctx context.Context) (*oidc.Provider, error) {
	return oidc.NewProvider(ctx, "https://access.line.me")
}

func lineConfig(provider *oidc.Provider) oauth2.Config {
	lineConfig := oauth2.Config{
		ClientID:     os.Getenv("LINE_CHANNEL_ID"),
		ClientSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		RedirectURL:  os.Getenv("LINE_REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return lineConfig
}

func AuthLineRedirect(w http.ResponseWriter, r *http.Request) {
	provider, err := lineProvider(r.Context())
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}

	lineConfig := lineConfig(provider)

	state, err := utils.RandString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	nonce, err := utils.RandString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	utils.SetCallbackCookie(w, r, "state", state)
	utils.SetCallbackCookie(w, r, "nonce", nonce)

	fmt.Println(lineConfig.AuthCodeURL(state, oidc.Nonce(nonce)))

	http.Redirect(w, r, lineConfig.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func AuthLineCallback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("state")
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		Value:    "",
		Path:     "/",
		Name:     "state",
	})
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	provider, err := lineProvider(r.Context())
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, "Internal errro", http.StatusInternalServerError)
	}

	config := lineConfig(provider)

	oauth2Token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	nonce, err := r.Cookie("nonce")
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		Value:    "",
		Path:     "/",
		Name:     "nonce",
	})
	if err != nil {
		http.Error(w, "nonce not found", http.StatusBadRequest)
		return
	}

	oauth2Token.AccessToken = "*REDACTED*"

	endPoint := "https://api.line.me/oauth2/v2.1/verify"
	values := url.Values{}
	values.Add("id_token", rawIDToken)
	values.Add("client_id", os.Getenv("LINE_CHANNEL_ID"))
	values.Add("nonce", nonce.Value)
	res, err := http.PostForm(endPoint, values)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	claims := authentication.Claims{
		Provider: "line",
	}
	if err := decoder.Decode(&claims); err != nil {
		return
	}

	dbPool := db.DatabaseConnection()
	isAuth, err := claims.CheckAuthAlreadyExists(r.Context(), dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isAuth {
		userID, err := user.GetUserIDByProviderAndUID(r.Context(), dbPool, "line", claims.ID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jwt, _ := user.CreateToken(userID)
		auth.SetAuthCookie(r.Context(), jwt)
		http.Redirect(w, r, os.Getenv("CLIENT_BASE_URL"), http.StatusPermanentRedirect)
	} else {
		userID, err := claims.CreateFrom(r.Context(), dbPool)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jwt, _ := user.CreateToken(userID)
		auth.SetAuthCookie(r.Context(), jwt)
		http.Redirect(w, r, os.Getenv("CLIENT_BASE_URL"), http.StatusPermanentRedirect)
	}
}
