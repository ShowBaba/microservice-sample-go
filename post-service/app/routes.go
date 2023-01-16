package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/microservice-sample-go/post-service/data"
	"github.com/showbaba/microservice-sample-go/shared"
)

var (
	models      *data.Models
	messageChan *amqp.Channel
	ctx         = context.Background()
)

func (a *App) Initialize(channel *amqp.Channel, dbModels *data.Models) {
	a.Router = mux.NewRouter()
	a.setRouters()
	models = dbModels
	messageChan = channel

	// declare the notifications exchange
	err := messageChan.ExchangeDeclare(shared.NOTIFICATION_TOPIC, "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) setRouters() {
	a.Router.Use(ValidateGatewayToken())
	a.Get("/post/ping", Ping)
	a.Post("/post/create", ValidateAuthToken(), CreatePost)

}

// handler method
func (a *App) Post(path string, mw func(http.Handler) http.Handler, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.Handle(path, mw(http.HandlerFunc(f)).(http.HandlerFunc)).Methods("Post")
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
