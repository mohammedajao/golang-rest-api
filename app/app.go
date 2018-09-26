package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/mohammedajao/rest-api/app/config"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type function func(w http.ResponseWriter, r *http.Request)

func (a *App) Initialize() {
	a.DB = config.InitDatabase()
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {
	a.Router.HandleFunc("/", index)
	http.Handle("/", a.Router)

	srv := &http.Server{
		Addr: "0.0.0.0:" + addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router, // Pass our instance of gorilla/mux in.
	}

	fmt.Printf("App is running on port: %s\n", addr)
	log.Fatal(srv.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Test")
}
