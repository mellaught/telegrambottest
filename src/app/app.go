package app

import (
	"database/sql"
	"log"
	"net/http"

	stct "github.com/mrKitikat/telegrambottest/src/app/bipdev/structs"
	"github.com/mrKitikat/telegrambottest/src/app/bot"

	"github.com/gorilla/mux"
)

// App is main app for ExchangeBot.
type App struct {
	Router *mux.Router
	Bot    *bot.Bot
}

// InitService is initializes the app.
func NewApp(conf *stct.Config, dbsql *sql.DB) *App {

	a := App{
		Router: mux.NewRouter(),
		Bot:    &bot.Bot{},
	}

	a.Router = mux.NewRouter()
	a.Bot = bot.InitBot(conf, dbsql)
	a.setRouters()

	return &a
}

// // Get wraps the router for GET method
// func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
// 	a.Router.HandleFunc(path, f).Methods("GET")
// }

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) setRouters() {

	// Routing for handling the Update user's loots.
	a.Post("/UpdateLoot", a.Bot.UpdateLoots)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
