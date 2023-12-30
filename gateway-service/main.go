package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/microservice-sample-go/gateway-service/app"
)

func main() {
	// setup rabbitmq
	connection, err := amqp.Dial(app.GetConfig().RabbitmqServerURL)
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
	app_.Initialize(messageChan)
	log.Printf("starting server on port: %s", port)
	app_.Run(port)
}
