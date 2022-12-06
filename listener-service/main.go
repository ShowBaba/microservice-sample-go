package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/microservice-sample-go/shared"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	channel *amqp.Channel
	ctx     = context.Background()
)

func main() {
	connection, err := amqp.Dial(shared.RABBITMQ_SERVER_URL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err = connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	//declare queue
	_, err = channel.QueueDeclare(
		shared.LISTENER_SERVICE,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	// subscribe to service
	messages, err := channel.Consume(
		shared.LISTENER_SERVICE,
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
			var payload shared.ListenerServicePayload
			if err = json.Unmarshal(message.Body, &payload); err != nil {
				panic(err)
			}
			handlePayload(payload)
		}
	}()
	log.Println("waiting for message")
	<-forever
}

func handlePayload(payload shared.ListenerServicePayload) {
	switch payload.ServiceName {
	case shared.NOTIFICATION_SERVICE:
		b, err := json.Marshal(payload.Data)
		if err != nil {
			log.Fatal(err)
		}
		message := amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(b),
		}
		if err = channel.PublishWithContext(
			ctx,
			"",
			shared.NOTIFICATION_SERVICE,
			false,
			false,
			message,
		); err != nil {
			log.Fatal(err)
		}
	case shared.LOGGER_SERVICE:
		b, err := json.Marshal(payload.Data)
		if err != nil {
			log.Fatal(err)
		}
		message := amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(b),
		}
		if err = channel.PublishWithContext(
			ctx,
			"",
			shared.LOGGER_SERVICE,
			false,
			false,
			message,
		); err != nil {
			log.Fatal(err)
		}
	default:
		log.Printf("unknown service name; %s", payload.ServiceName)
	}
}
