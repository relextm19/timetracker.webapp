package app

import (
	"context"
	"net/http"
	"slices"
	"strings"
)

var publicPaths = []string{"/login", "/register"}

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

// we need to declare a custom type for context to avoid colisions
type ctxKey string

const ctxKeyToken ctxKey = "authToken"

func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(publicPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		token := GetAuthTokenFromRequest(r)
		if token == "" {
			a.Logger.Warn("authentication failed: missing token", "path", r.URL.Path, "ip", r.RemoteAddr)
			a.RespondWithError(w, http.StatusUnauthorized, RespUnauthorized)
			return
		}

		valid, err := a.Store.CheckTokenValid(token)
		if err != nil {
			a.Logger.Error("failed to validate token in database", "error", err, "path", r.URL.Path)
			a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
			return
		}

		if valid {
			// save the token for use in different handlers
			ctx := context.WithValue(r.Context(), ctxKeyToken, token)
			reqWithContext := r.WithContext(ctx)
			next.ServeHTTP(w, reqWithContext)
			return
		} else {
			a.Logger.Warn("authentication failed: invalid token", "path", r.URL.Path, "ip", r.RemoteAddr)
			a.RespondWithError(w, http.StatusUnauthorized, RespUnauthorized)
			return
		}
	})
}
