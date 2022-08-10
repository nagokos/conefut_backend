package oauth

import (
	"context"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/cookie"
	"github.com/nagokos/connefut_backend/graph/jwt"
	"github.com/nagokos/connefut_backend/graph/models/authentication"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
	"golang.org/x/oauth2"
)

func googleProvider(ctx context.Context) (*oidc.Provider, error) {
	return oidc.NewProvider(ctx, "https://accounts.google.com")
}

func googleConfig(provider *oidc.Provider) oauth2.Config {
	googleConfig := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return googleConfig
}

func AuthGoogleRedirect(w http.ResponseWriter, r *http.Request) {
	provider, err := googleProvider(r.Context())
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}

	googleConfig := googleConfig(provider)

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

	http.Redirect(w, r, googleConfig.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func AuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
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
		logger.NewLogger().Error(err.Error())
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	provider, err := googleProvider(r.Context())
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, "Internal errro", http.StatusInternalServerError)
	}

	config := googleConfig(provider)
	oauth2Token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	}
	verifier := provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInsufficientStorage)
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
	if idToken.Nonce != nonce.Value {
		logger.NewLogger().Error(err.Error())
		http.Error(w, "nonce did not match", http.StatusBadRequest)
		return
	}

	oauth2Token.AccessToken = "*REDACTED*"

	claims := authentication.Claims{
		Provider: "google",
	}

	if err := idToken.Claims(&claims); err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbPool := db.DatabaseConnection()
	defer dbPool.Close()
	isAuth, err := claims.CheckAuthAlreadyExists(r.Context(), dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isAuth {
		userID, err := user.GetUserIDByProviderAndUID(r.Context(), dbPool, "google", claims.ID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jwt, _ := jwt.GenerateToken(userID)
		cookie.SetAuthCookie(r.Context(), jwt)
		http.Redirect(w, r, os.Getenv("CLIENT_BASE_URL"), http.StatusPermanentRedirect)
	} else {
		userID, err := claims.CreateFrom(r.Context(), dbPool)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jwt, _ := jwt.GenerateToken(userID)
		cookie.SetAuthCookie(r.Context(), jwt)
		http.Redirect(w, r, os.Getenv("CLIENT_BASE_URL"), http.StatusPermanentRedirect)
	}
}
