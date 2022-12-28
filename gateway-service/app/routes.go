package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	messageChan *amqp.Channel
	ctx         = context.Background()
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	// a.Router.Use(handler())
	a.Any("/{[^/]+}", Handler)
	a.Any("/{[^/]+}/{[^/]+}", Handler)
	a.Any("/{[^/]+}/{[^/]+}/{[^/]+}", Handler)
	a.Get("/gateway/ping", Ping)
}

func (a *App) Any(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Get", "Post", "Patch", "Put")
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Get")
}

func (a *App) Run(port string) {
	log.Fatal(
		http.ListenAndServe(
			port,
			handlers.CORS(
				handlers.AllowCredentials(),
				handlers.AllowedMethods([]string{"POST"}),
				handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
				handlers.MaxAge(3600),
			)(a.Router),
		),
	)
}
