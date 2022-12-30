package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	app_ "github.com/showbaba/microservice-sample-go/auth-service/app"
	"github.com/showbaba/microservice-sample-go/auth-service/data"
	"github.com/showbaba/microservice-sample-go/shared"
)

func main() {
	// setup rabbitmq
	connection, err := amqp.Dial(app_.GetConfig().RabbitmqServerURL)
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
	app := app_.App{}
	port := app_.GetConfig().Port
	app.Initialize(messageChan, &models)
	log.Printf("starting server on port: %s", port)
	app.Run(port)
}
