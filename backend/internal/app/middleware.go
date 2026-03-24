package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/relextm19/tracker.nvim/internal/helpers"
)

type RouteConfig struct {
	IsPublic   bool
	AllowKey   bool
	AllowToken bool
}

var routes = map[string]map[string]RouteConfig{
	"/login": {
		http.MethodGet:  {IsPublic: true},
		http.MethodPost: {IsPublic: true},
	},
	"/register": {
		http.MethodGet:  {IsPublic: true},
		http.MethodPost: {IsPublic: true},
	},
	"/checkAuth": {
		http.MethodGet:  {IsPublic: true},
		http.MethodPost: {IsPublic: true},
	},
	"/sessions": {
		http.MethodPost: {AllowKey: true},
		http.MethodGet:  {AllowKey: true, AllowToken: true}, // Both allowed so the nvim display can work
	},
}

// GetAuthTokenFromRequest since we have both browser and other clients making request we have to check for both cookies and headers
func GetAuthTokenFromRequest(r *http.Request) string {
	// TODO: i think we dont need to check for body anymore cuz only the browser auths by token
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
	ctxUserID ctxKey = "userID"
)

func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeConfig, ok := routes[r.URL.Path][r.Method]
		if !ok {
			a.Logger.Warn("Accesing non existent path")
			a.RespondWithError(w, http.StatusNotFound, RespNotFound)
			return
		}
		if routeConfig.IsPublic {
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

		isValid := false
		var dbErr error
		var userID string

		if routeConfig.AllowKey {
			apiKey := GetAPIKeyFromRequest(r)
			if apiKey != "" {
				hash, err := helpers.GetHashFromUUID([]byte(apiKey))
				if err == nil {
					userID, dbErr = a.Store.GetUserIDForKeyHash(hash)
					isValid = userID != ""
				}
			}
		}
		if !isValid && routeConfig.AllowToken {
			token := GetAuthTokenFromRequest(r)
			if token != "" {
				userID, dbErr = a.Store.GetUserIDForToken(token)
				isValid = userID != ""
			}
		}

		if dbErr != nil {
			failInternal("failed to validate credential in database", dbErr)
			return
		}

		if !isValid {
			failAuth("invalid credentials")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
