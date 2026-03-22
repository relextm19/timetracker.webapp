package app

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/relextm19/tracker.nvim/internal/helpers"
)

var (
	publicPaths    = []string{"/login", "/register"}
	keyAuthedPaths = []string{"/sessions"}
)

// GetAuthTokenFromRequest since we have both browser and other clients making request we have to check for both cookies and headers
func GetAuthTokenFromRequest(r *http.Request) string {
	h := r.Header
	if after, ok := strings.CutPrefix(h.Get("Authorization"), "Bearer"); ok {
		return strings.TrimSpace(after)
	}
	token, err := r.Cookie("token")
	if err == nil {
		return strings.TrimSpace(token.Value)
	}

	return ""
}

func GetAPIKeyFromRequest(r *http.Request) string {
	h := r.Header
	if after, ok := strings.CutPrefix(h.Get("Authorization"), "Bearer"); ok {
		return strings.TrimSpace(after)
	}

	return ""
}

// we need to declare a custom type for context to avoid colisions
type ctxKey string

const (
	ctxKeyToken  ctxKey = "authToken"
	ctxKeyAPIKey ctxKey = "APIKey"
)

func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(publicPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		failAuth := func(reason string) {
			a.Logger.Warn("authentication failed: "+reason, "path", r.URL.Path, "ip", r.RemoteAddr)
			a.RespondWithError(w, http.StatusUnauthorized, RespUnauthorized)
		}

		failInternal := func(msg string, err error) {
			a.Logger.Error(msg, "error", err, "path", r.URL.Path)
			a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		}

		var isValid bool
		var dbErr error
		var ctxKey any
		var ctxValue string

		if slices.Contains(keyAuthedPaths, r.URL.Path) {
			apiKey := GetAPIKeyFromRequest(r)
			if apiKey == "" {
				failAuth("missing api key")
				return
			}

			hash, err := helpers.GetHashFromUUID([]byte(apiKey))
			if err != nil {
				failInternal("error getting key hash", err)
				return
			}

			isValid, dbErr = a.Store.CheckKeyHashValid(hash)
			ctxKey, ctxValue = ctxKeyAPIKey, hash
		} else {
			token := GetAuthTokenFromRequest(r)
			if token == "" {
				failAuth("missing token")
				return
			}

			isValid, dbErr = a.Store.CheckTokenValid(token)
			ctxKey, ctxValue = ctxKeyToken, token
		}

		if dbErr != nil {
			failInternal("failed to validate credential in database", dbErr)
			return
		}

		if !isValid {
			failAuth("invalid credential")
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey, ctxValue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
