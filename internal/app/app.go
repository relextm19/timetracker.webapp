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

func (a *App) HandleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			a.Logger.Error(err.Error())
			return
		}

		fmt.Println(string(body))
		session := sessions.NewSession()
		json.Unmarshal(body, session)
		if err = session.IsValid(); err != nil {
			a.Logger.Error(err.Error())
			return
		}

		if err = a.Store.InsertSession(session); err != nil {
			a.Logger.Error(err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (a *App) HandleRegister(w http.ResponseWriter, r *http.Request) {
}
