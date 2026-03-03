package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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

func (a *App) respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	a.Logger.Error(err.Error())
	http.Error(w, msg, code)
}

func (a *App) SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.respondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed, nil)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespBadRequest, err)
		return
	}

	fmt.Println(string(body))
	session := sessions.NewSession()

	err = json.Unmarshal(body, session)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespInvalidJSON, err)
		return
	}

	if err = session.Valid(); err != nil {
		a.respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err = a.Store.InsertSession(session); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.respondWithError(w, http.StatusMethodNotAllowed, RespMethodNotAllowed, nil)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespBadRequest, err)
		return
	}

	cub := users.NewClientUserBody()
	err = json.Unmarshal(body, cub)
	if err != nil {
		a.respondWithError(w, http.StatusBadRequest, RespInvalidJSON, err)
		return
	}

	user, err := users.NewUser(cub)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError, err)
		return
	}

	if err = user.Valid(); err != nil {
		a.respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err = a.Store.InsertUser(user); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, RespInternalError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
