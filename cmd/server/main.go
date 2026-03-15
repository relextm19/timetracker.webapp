package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/relextm19/tracker.nvim/internal/app"
	"github.com/rs/cors"
)

func main() {
	dbPath := flag.String("db path", "./db/database.db", "path to the db file")
	flag.Parse()

	app := app.NewApp(dbPath)
	defer app.Store.DB.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/sessions", app.SessionHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := app.AuthMiddleware(c.Handler(mux))

	log.Panic(http.ListenAndServe(":42069", handler))
}
