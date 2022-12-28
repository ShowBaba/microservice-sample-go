package main

import (
	"log"

	"github.com/microservice-sample-go/gateway-service/app"
	"github.com/microservice-sample-go/shared"
	amqp "github.com/rabbitmq/amqp091-go"
)



func main() {
	// setup rabbitmq
	connection, err := amqp.Dial(shared.RABBITMQ_SERVER_URL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	messageChan, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer messageChan.Close()
	app_ := app.App{}
	port := app.GetConfig().Port
	app_.Initialize()
	log.Printf("starting server on port: %s", port)
	app_.Run(port)
}
