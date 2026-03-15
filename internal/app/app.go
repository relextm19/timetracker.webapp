package app

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	database "github.com/relextm19/tracker.nvim/internal/db"
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
