package app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/microservice-sample-go/shared"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func LogRequest(source, data string) error {
	var entry struct {
		Source string `json:"source"`
		Data   string `json:"data"`
	}
	entry.Source = source
	entry.Data = data
	messageData := shared.ListenerServicePayload{
		ServiceName: shared.LOGGER_SERVICE,
		Data:        entry,
	}
	b, err := json.Marshal(messageData)
	if err != nil {
		fmt.Println("err marshal - ", err)
		log.Fatal(err)
	}
	// create a message to publish
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(b),
	}
	// publish message to queue
	if err = messageChan.PublishWithContext(
		ctx,
		"",
		shared.LISTENER_SERVICE,
		false,
		false,
		message,
	); err != nil {
		fmt.Println("err publishing message ", err)
		return err
	}
	return nil
}
