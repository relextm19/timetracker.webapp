package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/relextm19/tracker.nvim/internal/sessions"
	"github.com/relextm19/tracker.nvim/internal/users"
)

const (
	RespMethodNotAllowed = "Method not allowed"
	RespBadRequest       = "Bad request"
	RespInvalidJSON      = "Invalid JSON"
	RespInternalError    = "Internal server error"
	RespUnauthorized     = "Incorrect credentials"
)

var ErrSessionInvalid = fmt.Errorf("invalid session")

func (a *App) RespondWithError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}

func setAuthCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func (a *App) SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.CreateSession(w, r)
	case http.MethodGet:
		a.GetUserData(w, r)
	default:
		a.RespondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
	}
}

func (a *App) CreateSession(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error("failed to read request body", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}
	defer r.Body.Close()

	session := sessions.NewSession()

	err = json.Unmarshal(body, session)
	if err != nil {
		a.Logger.Error("failed to unmarshal session", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = session.Valid(); err != nil {
		a.Logger.Error("invalid session", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, _ := r.Context().Value(ctxKeyToken).(string) // this should never error

	if err = a.Store.InsertSession(session, token); err != nil {
		a.Logger.Error("failed to insert session", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) GetUserData(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(ctxKeyToken).(string) // this should never error

	data, err := a.Store.GetDataForToken(token)
	if err != nil {
		a.Logger.Error("failed to get data for token", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.RespondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error("failed to read request body", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	cub := users.NewClientUserBody()
	err = json.Unmarshal(body, cub)
	if err != nil {
		a.Logger.Error("failed to unmarshal user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = cub.Valid(); err != nil {
		a.Logger.Error("invalid user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := users.NewUser(cub)
	if err != nil {
		a.Logger.Error("failed to create user", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	if err = a.Store.InsertUser(user); err != nil {
		a.Logger.Error("failed to insert user", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	setAuthCookie(w, user.Token.String())
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(struct {
		Token uuid.UUID `json:"token"`
	}{
		Token: user.Token,
	})
}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.RespondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error("failed to read request body", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	cub := users.NewClientUserBody()
	err = json.Unmarshal(body, cub)
	if err != nil {
		a.Logger.Error("failed to unmarshal user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = cub.Valid(); err != nil {
		a.Logger.Error("invalid user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if ok := a.Store.CheckLoginAttempt(cub); !ok {
		a.Logger.Error("login attempt failed", "email", cub.Email)
		a.RespondWithError(w, http.StatusUnauthorized, RespUnauthorized)
		return
	}

	token, err := a.Store.GetUserToken(cub.Email)
	if err != nil {
		a.Logger.Error("failed to fetch user token", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	setAuthCookie(w, token.String())
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(struct {
		Token uuid.UUID `json:"token"`
	}{
		Token: token,
	})
}
