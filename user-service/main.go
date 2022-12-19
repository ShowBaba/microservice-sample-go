package main

import (
	"log"

	"github.com/microservice-sample-go/shared"
	app_ "github.com/microservice-sample-go/user-service/app"
	"github.com/microservice-sample-go/user-service/data"
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
		app_.GetConfig().DbHost,
		app_.GetConfig().DbUser,
		app_.GetConfig().DbPassword,
		app_.GetConfig().DbName,
		app_.GetConfig().DbPort,
	)
	defer dbConn.Close()
	models := data.New(dbConn)
	data.Migrate()
	app := app_.App{}
	port := app_.GetConfig().Port
	app.Initialize(messageChan, &models)
	log.Printf("starting server on port: %s", port)
	app.Run(port)
}
