package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	apikeys "github.com/relextm19/tracker.nvim/internal/apiKeys"
	database "github.com/relextm19/tracker.nvim/internal/db"
	"github.com/relextm19/tracker.nvim/internal/sessions"
	"github.com/relextm19/tracker.nvim/internal/users"
)

const (
	RespMethodNotAllowed = "Method not allowed"
	RespNotFound         = "Path not found"
	RespBadRequest       = "Bad request"
	RespInvalidJSON      = "Invalid JSON"
	RespInternalError    = "Internal server error"
	RespUnauthorized     = "Incorrect credentials"
)

var ErrSessionInvalid = fmt.Errorf("invalid session")

// if the function fails it means that the server is setup Incorrectly in some way and shouldnt continue so panic
func getTokenFromContext(r *http.Request) string {
	token, ok := r.Context().Value(ctxKeyToken).(string)
	if !ok {
		panic("failed to get token from context - auth middleware may not be set up correctly")
	}
	return token
}

// if the function fails it means that the server is setup Incorrectly in some way and shouldnt continue so panic
func getAPIKeyFromContext(r *http.Request) string {
	key, ok := r.Context().Value(ctxKeyAPIKey).(string)
	if !ok {
		panic("failed to get api key from context - auth middleware may not be set up correctly")
	}
	return key
}

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

	apiKey := getAPIKeyFromContext(r)

	if err = a.Store.InsertSession(session, apiKey); err != nil {
		a.Logger.Error("failed to insert session", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) GetUserData(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromContext(r)

	data, err := a.Store.GetSessionDataForToken(token)
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

	rub := users.NewRequestUserBody()
	err = json.Unmarshal(body, rub)
	if err != nil {
		a.Logger.Error("failed to unmarshal user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = rub.Valid(); err != nil {
		a.Logger.Error("invalid user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := users.NewUser(rub)
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

	rub := users.NewRequestUserBody()
	err = json.Unmarshal(body, rub)
	if err != nil {
		a.Logger.Error("failed to unmarshal user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = rub.Valid(); err != nil {
		a.Logger.Error("invalid user", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if ok := a.Store.CheckLoginAttempt(rub); !ok {
		a.Logger.Error("login attempt failed", "email", rub.Email)
		a.RespondWithError(w, http.StatusUnauthorized, RespUnauthorized)
		return
	}

	token, err := a.Store.GetUserToken(rub.Email)
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

func (a *App) AddAPIKey(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error("failed to read request body", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	rak := apikeys.NewRequestAPIKey()
	err = json.Unmarshal(body, rak)
	if err != nil {
		a.Logger.Error("failed to unmarshal api key", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	err = rak.Valid()
	if err != nil {
		a.Logger.Error("invalid client api key", "error", err)
		a.RespondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	ak, err := apikeys.NewAPIKey(rak)
	if err != nil {
		a.Logger.Error("failed to create api key", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	token := getTokenFromContext(r)

	id, createdAt, err := a.Store.InsertAPIKey(token, ak)
	if err != nil {
		if errors.Is(err, database.ErrNoRowsAffected) {
			a.Logger.Warn("failed insert api key", "error", err)
			a.RespondWithError(w, http.StatusNotFound, RespBadRequest)
			return
		}
		a.Logger.Error("failed insert api key", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}
	ak.ID = id
	ak.CreatedAt = createdAt

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ak)
}

func (a *App) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	keyID := r.PathValue("id")
	if keyID == "" {
		a.Logger.Error("failed to retrieve api key from path", "error", "missing path parameter")
		a.RespondWithError(w, http.StatusBadRequest, "missing key id")
		return
	}

	token := getTokenFromContext(r)
	err := a.Store.DeleteAPIKey(keyID, token)
	if err != nil {
		if errors.Is(err, database.ErrNoRowsAffected) {
			a.Logger.Warn("failed to delete api key", "error", err)
			a.RespondWithError(w, http.StatusNotFound, RespBadRequest)
			return
		}
		a.Logger.Error("failed to delete api key", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) GetAPIKeys(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromContext(r)
	keys, err := a.Store.GetAPIKeys(token)
	if err != nil {
		a.Logger.Error("failed to fetch api keys", "error", err)
		a.RespondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}

func (a *App) APIKeysHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.AddAPIKey(w, r)
	case http.MethodGet:
		a.GetAPIKeys(w, r)
	default:
		a.RespondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
	}
}
