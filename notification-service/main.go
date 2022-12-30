package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/microservice-sample-go/notification-service/app"
	"github.com/showbaba/microservice-sample-go/shared"
)

var (
	channel *amqp.Channel
	ctx     = context.Background()
)

func main() {
	connection, err := amqp.Dial(app.GetConfig().RabbitmqServerURL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	channel, err = connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(shared.NOTIFICATION_TOPIC, "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	emailQueue, err := channel.QueueDeclare("email", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}
	smsQueue, err := channel.QueueDeclare("sms", false, false, true, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(emailQueue.Name, "email.*", shared.NOTIFICATION_TOPIC, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.QueueBind(smsQueue.Name, "sms.*", shared.NOTIFICATION_TOPIC, false, nil)
	if err != nil {
		panic(err)
	}
	// TODO: remove autoAck and setup manual acknowledgement
	emailMsgs, err := channel.Consume(emailQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	smsMsgs, err := channel.Consume(smsQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening for notifications...")
	for {
		select {
		case emailMsg := <-emailMsgs:
			// Send an acknowledgement message
			// err = channel.Ack(emailMsg.DeliveryTag, false)
			// if err != nil {
			// 	fmt.Println(err)
			// 	if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("err; %v", err)); err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	log.Println(err)
			// }
			var payload app.EmailMsgPayload
			err := json.Unmarshal(emailMsg.Body, &payload)
			if err != nil {
				if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("err; %v", err)); err != nil {
					log.Fatal(err)
				}
				log.Fatal(err)
			}
			if err := app.HandleEmailMsg(ctx, channel, payload); err != nil {
				if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("err; %v", err)); err != nil {
					log.Fatal(err)
				}
				log.Fatal(err)
			}

		case smsMsg := <-smsMsgs:
			fmt.Printf("Received sms notification: %s\n", smsMsg.Body)
			// err = channel.Ack(smsMsg.DeliveryTag, false)
			// if err != nil {
			// 	if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("err; %v", err)); err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	log.Println(err)
			// }
		}
	}
}
