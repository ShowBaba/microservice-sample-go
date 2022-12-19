package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/microservice-sample-go/auth-service/data"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	models      data.Models
	messageChan *amqp.Channel
	ctx         = context.Background()
)

func (a *App) Initialize(channel *amqp.Channel, dbModels *data.Models) {
	a.Router = mux.NewRouter()
	a.setRouters()
	models = *dbModels
	messageChan = channel
}

func (a *App) setRouters() {
	a.Router.Use(ValidateGatewayToken())
	a.Post("/auth/login", Login)
	a.Get("/auth/ping", Ping)
}

// handler method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Post")
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Get")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Put")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Delete")
}

// run
func (a *App) Run(host string) {
	// CORS
	log.Fatal(
		http.ListenAndServe(
			host,
			handlers.CORS(
				handlers.AllowCredentials(),
				handlers.AllowedMethods([]string{"POST", "GET", "PUT", "OPTIONS"}),
				handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
				handlers.MaxAge(3600),
			)(a.Router),
		),
	)
}
