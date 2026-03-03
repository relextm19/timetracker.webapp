package main

import (
	"flag"
	"net/http"

	"github.com/relextm19/tracker.nvim/internal/app"
)

func main() {
	dbPath := flag.String("db path", "./db/database.db", "path to the db file")
	flag.Parse()

	app := app.NewApp(dbPath)
	defer app.Store.DB.Close()

	http.HandleFunc("/sessions", app.SessionHandler)
	http.HandleFunc("/users", app.RegisterHandler)
	app.Logger.Info(http.ListenAndServe(":42069", nil).Error())
}
