package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/microservice-sample-go/notification-service/app"
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
		shared.NOTIFICATION_SERVICE,
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
		shared.NOTIFICATION_SERVICE,
		"",
		false, // set auto ack to false so that messages can be manually ack when service is done processing message
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
			var payload Payload
			if err = json.Unmarshal(message.Body, &payload); err != nil {
				panic(err)
			}
			if err := handleMessage(payload); err != nil {
				if err := message.Ack(true); err != nil {
					log.Fatal(err)
				}
				if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("err; %v", err)); err != nil {
					log.Fatal(err)
				}
				log.Fatal(err)
			}
			if err := message.Ack(false); err != nil {
				log.Fatal(err)
			}
		}
	}()
	log.Println("waiting for message")
	<-forever
}

type Payload struct {
	Channel string
	Data    struct {
		To      []string
		Subject string
		Body    string
	}
}

func handleMessage(payload Payload) error {
	switch payload.Channel {
	case shared.MAIL:
		mail := app.Mail{
			Sender:  shared.MAIL_USERNAME,
			Subject: payload.Data.Subject,
			To:      payload.Data.To,
			Body:    payload.Data.Body,
		}
		if err := app.SendEmail(mail); err != nil {
			if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("error while sending mail; err: %v", err)); err != nil {
				return err
			}
			return fmt.Errorf("error while sending mail; err: %v", err)
		}
		if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("sent email; payload: %v", payload)); err != nil {
			return err
		}
	case shared.SMS:
	default:
		if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("%v", payload)); err != nil {
			return err
		}
	}
	return nil
}
