package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	database "github.com/relextm19/tracker.nvim/internal/db"
	sessions "github.com/relextm19/tracker.nvim/internal/sessions"
	"github.com/relextm19/tracker.nvim/internal/users"
)

const (
	RespMethodNotAllowed = "Method not allowed"
	RespBadRequest       = "Bad request"
	RespInvalidJSON      = "Invalid JSON"
	RespInternalError    = "Internal server error"
	RespUnauthorized     = "Incorrect credentials"
)

type App struct {
	Store  *database.Store
	Logger slog.Logger
}

func NewApp(dbPath *string) *App {
	a := &App{}

	a.Logger = *slog.Default()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		panic(err)
	}
	a.Store = database.NewStore(db)

	return a
}

var ErrSessionInvalid = fmt.Errorf("invalid session")

func (a *App) respondWithError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}

func (a *App) SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.respondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	session := sessions.NewSession()

	err = json.Unmarshal(body, session)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = session.Valid(); err != nil {
		a.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = a.Store.InsertSession(session); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.respondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	cub := users.NewClientUserBody()
	err = json.Unmarshal(body, cub)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = cub.Valid(); err != nil {
		a.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := users.NewUser(cub)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	if err = a.Store.InsertUser(user); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(struct {
		Token uuid.UUID `json:"token"`
	}{
		Token: user.Token,
	})
}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.respondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespBadRequest)
		return
	}

	cub := users.NewClientUserBody()
	err = json.Unmarshal(body, cub)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespInvalidJSON)
		return
	}

	if err = cub.Valid(); err != nil {
		a.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if ok := a.Store.CheckLoginAttempt(cub); !ok {
		a.respondWithError(w, http.StatusUnauthorized, RespUnauthorized)
		return
	}

	token, err := a.Store.GetUserToken(cub.Email)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(struct {
		Token uuid.UUID `json:"token"`
	}{
		Token: token,
	})
}

func (a *App) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the actual handler (e.g., LoginHandler or RegisterHandler)
		next.ServeHTTP(w, r)

		// Log after the handler finishes
		a.Logger.Info("handled request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
