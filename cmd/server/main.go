package main

import (
	"net/http"

	"github.com/relextm19/tracker.nvim/internal/app"
)

func main() {
	app := app.NewApp()
	http.HandleFunc("/", app.HandleHome)
	http.ListenAndServe(":2137", nil)
}
