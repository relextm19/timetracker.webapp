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
	mux := http.NewServeMux()
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/session", app.SessionHandler)

	// loggedMux := app.LoggingMiddleware(mux)

	http.ListenAndServe(":42069", mux)
}
