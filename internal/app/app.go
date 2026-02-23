package app

import (
	"database/sql"
	"fmt"
	"net/http"
)

type App struct {
	SecretVal int
	DB        sql.DB
}

func NewApp() *App {
	a := &App{}
	a.SecretVal = 67
	return a
}

func (a *App) HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(a.SecretVal)
}
