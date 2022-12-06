package main

import (
	"log"

	"github.com/microservice-sample-go/auth-service/app"
	"github.com/microservice-sample-go/shared"
	"github.com/microservice-sample-go/auth-service/data"
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
	dbConn := shared.ConnectToSQLDB(
		app.GetConfig().DbHost,
		app.GetConfig().DbUser,
		app.GetConfig().DbPassword,
		app.GetConfig().DbName,
		app.GetConfig().DbPort,
	)
	defer dbConn.Close()
	models := data.New(dbConn)
	app := app.App{}
	port := ":3000"
	app.Initialize(messageChan, &models)
	log.Printf("starting server on port: %s", port)
	app.Run(port)
}
