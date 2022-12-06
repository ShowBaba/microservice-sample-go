package main

import (
	"encoding/json"
	"log"

	"github.com/microservice-sample-go/logger-service/app"
	"github.com/microservice-sample-go/logger-service/data"
	"github.com/microservice-sample-go/shared"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	models *data.Models
)

func main() {
	client, ctx, cancel, err := data.ConnectToDB(app.GetConfig().MongoURI)
	data.PingDB(ctx, client)
	mongoInst := data.NewMongoClient(client)
	if err != nil {
		panic(err)
	}
	defer data.CloseDBConnection(client, ctx, cancel)
	models = &mongoInst
	connection, err := amqp.Dial(shared.RABBITMQ_SERVER_URL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	//declare queue
	_, err = channel.QueueDeclare(
		shared.LOGGER_SERVICE,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	messages, err := channel.Consume(
		shared.LOGGER_SERVICE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("error subscribing to message - %v", err)
	}
	forever := make(chan bool)
	go func() {
		for message := range messages {
			log.Printf(" > Received message: %s\n", message.Body)
			var payload app.LogPayload
			if err = json.Unmarshal(message.Body, &payload); err != nil {
				panic(err)
			}
			handlePayload(payload)
		}
	}()
	log.Println("waiting for message")
	<-forever
}


func handlePayload(payload app.LogPayload) {
	logData := data.Log{
		Data:   payload.Data,
		Source: payload.Source,
	}
	if err := models.Log.Insert(logData); err != nil {
		log.Fatal(err)
	}
}
